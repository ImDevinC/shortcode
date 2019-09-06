import os
import logging
import base64
import json
from http import HTTPStatus
from db import DB

logging.getLogger().setLevel(os.environ.get('LOG_LEVEL', 'INFO'))


def generateNewShortcode(event):
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
    if len(path) == 0:
        logging.debug('Create new shortlink with description')
        db = DB()
        err = db.insert_short_code(path, uri)
    else:
        code, err = createRandomShortlink(uri)

    if err:
        return client_error(HTTPStatus.INTERNAL_SERVER_ERROR)

    return link_created_response(code, uri)


def createRandomShortlink(uri):
    logging.debug('Create new random shortlink')
    db = DB()
    code = db.get_short_code_from_uri(uri)
    if code:
        return code, None
    for _ in range(0, 5):
        code = DB.generate_short_code(5)
        err = db.insert_short_code(code, uri)
        if err:
            continue
        break
    return code, err


def main(event, _):
    if event['httpMethod'] == 'POST':
        logging.info('POST event received')
        logging.debug(generateNewShortcode(event))
    elif event['httpMethod'] == 'GET':
        logging.info('GET event received')
    elif event['httpMethod'] == 'DELETE':
        logging.info('DELETE event received')
    else:
        logging.error('Unknown event type received')


def client_error(status_code):
    return {
        'statusCode': status_code,
        'body': HTTPStatus(status_code).phrase
    }


def link_created_response(shortcode, uri):
    return {
        'statusCode': HTTPStatus.CREATED,
        'body': json.dumps({'shortCode': shortcode, 'uri': uri})
    }


if __name__ == '__main__':
    sample_event = {
        'body': 'eyJ0ZXN0IjoiYm9keSJ9',
        'resource': '/{proxy+}',
        'path': '/',
        'httpMethod': 'POST',
        'isBase64Encoded': True,
        'queryStringParameters': {
            'foo': 'bar'
        },
        'multiValueQueryStringParameters': {
            'foo': [
                'bar'
            ]
        },
        'pathParameters': {
            'proxy': '/path/to/resource'
        },
        'stageVariables': {
            'baz': 'qux'
        },
        'body': 'eyJ1cmkiOiAiaHR0cHM6Ly9pbWRldmluYy5jb20ifQo=',
        'headers': {
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Encoding': 'gzip, deflate, sdch',
            'Accept-Language': 'en-US,en;q=0.8',
            'Cache-Control': 'max-age=0',
            'CloudFront-Forwarded-Proto': 'https',
            'CloudFront-Is-Desktop-Viewer': 'true',
            'CloudFront-Is-Mobile-Viewer': 'false',
            'CloudFront-Is-SmartTV-Viewer': 'false',
            'CloudFront-Is-Tablet-Viewer': 'false',
            'CloudFront-Viewer-Country': 'US',
            'Host': '1234567890.execute-api.us-east-1.amazonaws.com',
            'Upgrade-Insecure-Requests': '1',
            'User-Agent': 'Custom User Agent String',
            'Via': '1.1 08f323deadbeefa7af34d5feb414ce27.cloudfront.net (CloudFront)',
            'X-Amz-Cf-Id': 'cDehVQoZnx43VYQb9j2-nvCh-9z396Uhbp027Y2JvkCPNLmGJHqlaA==',
            'X-Forwarded-For': '127.0.0.1, 127.0.0.2',
            'X-Forwarded-Port': '443',
            'X-Forwarded-Proto': 'https'
        },
        'multiValueHeaders': {
            'Accept': [
                'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'
            ],
            'Accept-Encoding': [
                'gzip, deflate, sdch'
            ],
            'Accept-Language': [
                'en-US,en;q=0.8'
            ],
            'Cache-Control': [
                'max-age=0'
            ],
            'CloudFront-Forwarded-Proto': [
                'https'
            ],
            'CloudFront-Is-Desktop-Viewer': [
                'true'
            ],
            'CloudFront-Is-Mobile-Viewer': [
                'false'
            ],
            'CloudFront-Is-SmartTV-Viewer': [
                'false'
            ],
            'CloudFront-Is-Tablet-Viewer': [
                'false'
            ],
            'CloudFront-Viewer-Country': [
                'US'
            ],
            'Host': [
                '0123456789.execute-api.us-east-1.amazonaws.com'
            ],
            'Upgrade-Insecure-Requests': [
                '1'
            ],
            'User-Agent': [
                'Custom User Agent String'
            ],
            'Via': [
                '1.1 08f323deadbeefa7af34d5feb414ce27.cloudfront.net (CloudFront)'
            ],
            'X-Amz-Cf-Id': [
                'cDehVQoZnx43VYQb9j2-nvCh-9z396Uhbp027Y2JvkCPNLmGJHqlaA=='
            ],
            'X-Forwarded-For': [
                '127.0.0.1, 127.0.0.2'
            ],
            'X-Forwarded-Port': [
                '443'
            ],
            'X-Forwarded-Proto': [
                'https'
            ]
        },
        'requestContext': {
            'accountId': '123456789012',
            'resourceId': '123456',
            'stage': 'prod',
            'requestId': 'c6af9ac6-7b61-11e6-9a41-93e8deadbeef',
            'requestTime': '09/Apr/2015:12:34:56 +0000',
            'requestTimeEpoch': 1428582896000,
            'identity': {
                'cognitoIdentityPoolId': None,
                'accountId': None,
                'cognitoIdentityId': None,
                'caller': None,
                'accessKey': None,
                'sourceIp': '127.0.0.1',
                'cognitoAuthenticationType': None,
                'cognitoAuthenticationProvider': None,
                'userArn': None,
                'userAgent': 'Custom User Agent String',
                'user': None
            },
            'path': '/prod/path/to/resource',
            'resourcePath': '/{proxy+}',
            'httpMethod': 'POST',
            'apiId': '1234567890',
            'protocol': 'HTTP/1.1'
        }
    }
    main(sample_event, None)
