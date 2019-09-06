import os
from random import choice
import boto3
import logging


class DB:
    charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890'
    clean_charset = 'abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789'

    def __init__(self):
        self.client = boto3.client('dynamodb')

    @staticmethod
    def generate_short_code(size):
        return ''.join(choice(DB.clean_charset) for i in range(size))

    def get_short_code_from_uri(self, uri):
        try:
            response = self.client.query(
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
                IndexName='URIIndex',
                KeyConditionExpression='URI = :uri',
                ExpressionAttributeValues={
                    ':uri': {
                        'S': uri
                    }
                }
            )

            if not 'Items' in response or not response['Items']:
                return None, None

            return response['Items'], None
        except self.client.exceptions.ResourceNotFoundException:
            logging.debug('Existing shortcode for %s not found', uri)
            return None, None
        except Exception as ex:
            logging.error(ex)
            return None, ex

    def get_shortcode(self, code):
        try:
            response = self.client.query(
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
                KeyConditionExpression='shortcode = :shortcode',
                ExpressionAttributeValues={
                    ':shortcode': {
                        'S': code
                    }
                }
            )

            if not 'Items' in response or not response['Items']:
                return None, None

            return response['Items'][0], None
        except self.client.exceptions.ResourceNotFoundException:
            logging.debug('Existing shortcode %s not found', code)
            return None, None
        except Exception as ex:
            logging.error(ex)
            return None, ex

    def insert_short_code(self, code, uri, custom):
        try:
            self.client.put_item(
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
                Item={
                    'shortcode': {'S': code},
                    'URI': {'S': uri},
                    'custom': {'BOOL': custom}
                },
                ConditionExpression='attribute_not_exists(shortcode)'
            )
        except self.client.exceptions.ConditionalCheckFailedException as ex:
            return ex
        except Exception as ex:
            logging.error(ex)
            return ex
