import os
import logging
import base64
import json
from http import HTTPStatus
from db import DB

logging.getLogger().setLevel(os.environ.get('LOG_LEVEL', 'INFO').upper())


def generate_new_shortcode(event):
    body = json.loads(event['body'] if not event['isBase64Encoded']
                      else base64.b64decode(event['body']))
    if not 'uri' in body:
        return client_error(HTTPStatus.BAD_REQUEST)
    uri = body['uri']
    path = event['path']
    if path.startswith('/'):
        path = path[1:]
    if path.count('/') > 0:
        return client_error(HTTPStatus.BAD_REQUEST)
    if path:
        code, status = create_named_link(path, uri)
        if status:
            return client_error(status)
    else:
        code, err = create_random_shortlink(uri)
        if err:
            return client_error(HTTPStatus.INTERNAL_SERVER_ERROR)

    return link_created_response(code, uri)


def create_named_link(code, uri):
    database = DB()
    exists, _ = database.get_shortcode(code)
    if exists:
        return None, HTTPStatus.ALREADY_REPORTED

    err = database.insert_short_code(code, uri, True)
    if err:
        return None, HTTPStatus.INTERNAL_SERVER_ERROR

    return code, None


def create_random_shortlink(uri):
    logging.debug('Create new random shortlink for %s', uri)
    database = DB()
    codes, _ = database.get_short_code_from_uri(uri)
    for code in codes:
        if code and 'custom' in code and not code['custom']['BOOL']:
            logging.info('Found code %s', code['shortcode']['S'])
            return code['shortcode']['S'], None
    for _ in range(0, 5):
        code = DB.generate_short_code(5)
        logging.info('%s %s', code, uri)
        err = database.insert_short_code(code, uri, False)
        if err:
            continue
        break
    return code, err


def get_uri_from_shortcode(event):
    path = event['path']
    if path.startswith('/'):
        path = path[1:]
    if path.count('/') > 0 or not path:
        return return_redirect(os.environ['HOMEPAGE'])
    code, err = DB().get_shortcode(path)
    if err:
        return client_error(HTTPStatus.INTERNAL_SERVER_ERROR)
    return return_redirect(code['URI']['S'] if code else os.environ['HOMEPAGE'])


def main(event, _):
    logging.debug(event)
    response = client_error(HTTPStatus.METHOD_NOT_ALLOWED)
    if event['httpMethod'] == 'POST':
        logging.info('POST event received')
        response = generate_new_shortcode(event)
    elif event['httpMethod'] == 'GET':
        logging.info('GET event received')
        response = get_uri_from_shortcode(event)
    elif event['httpMethod'] == 'DELETE':
        logging.info('DELETE event received')
    else:
        logging.error('Unknown event type received')
    return response


def client_error(status_code):
    return {
        'statusCode': status_code,
        'statusDescription': HTTPStatus(status_code).phrase,
        'body': HTTPStatus(status_code).phrase,
        'headers': {
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Headers': '*',
        }
    }


def link_created_response(shortcode, uri):
    return {
        'statusCode': HTTPStatus.CREATED,
        'statusDescription': HTTPStatus(HTTPStatus.CREATED).phrase,
        'body': json.dumps({'shortCode': shortcode, 'uri': uri}),
        'headers': {
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Headers': '*',
        }
    }


def return_redirect(url):
    return {
        'statusCode': HTTPStatus.FOUND,
        'statusDescription': HTTPStatus(HTTPStatus.FOUND).phrase,
        'headers': {
            'Location': url,
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Headers': '*',
        }
    }


if __name__ == '__main__':
    sample_event = {
        'requestContext': {
            'elb': {
                'targetGroupArn': 'arn:aws:elasticloadbalancing:us-east-1:950747459100:targetgroup/shortcode-dev/feeafae259184b4d'
            }
        },
        'httpMethod': 'GET',
        'path': '/imdevinc',
        'queryStringParameters': {},
        'headers': {
            'accept': '*/*',
            'accept-encoding': 'gzip, deflate',
            'connection': 'keep-alive',
            'host': 'shortcode-dev-1630588222.us-east-1.elb.amazonaws.com',
            'user-agent': 'HTTPie/1.0.2',
            'x-amzn-trace-id': 'Root=1-5d755912-6ec8714001ab3380518a8fc0',
            'x-forwarded-for': '162.223.36.138',
            'x-forwarded-port': '80',
            'x-forwarded-proto': 'http'
        },
        'body': '',
        'isBase64Encoded': False
    }
    main(sample_event, None)
