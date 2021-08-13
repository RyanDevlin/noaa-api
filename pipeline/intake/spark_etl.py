import argparse
import os
from pyspark.sql import SparkSession
from pyspark import SparkContext, SparkConf

from intake.etl import run_etl


def spark_run_etl(source, output_path, aws_access_key_id='', aws_secret_access_key='', local=False):
    """
    Run ETL from Source and output to Parquet

    :params source (str) - source from intake/sources
    :params output_path (str) - output path to s3 or local file system
    :params aws_access_key_id (str) - defaults to None
    :params aws_secret_access_key (str) - defaults to None
    :params local (bool) - Run pipeline for AWS or Local (s3 or local file system)
    """
    # This is a work-around since we are running spark locally on EC2
    # If we were running on a hadoop cluster, we could bypass this...
    # Unfortunately, we are cheap and spend too much of our money on
    # NYC rent and street food...
    if not local:
        os.environ['PYSPARK_SUBMIT_ARGS'] = '--packages org.apache.hadoop:hadoop-common:3.0.0,org.apache.hadoop:hadoop-aws:3.0.0,org.apache.hadoop:hadoop-client:3.0.0 pyspark-shell'
        # Seems like we need to export env vars, too. Another
        # hacky workaround that will stick for now...
        os.environ['AWS_ACCESS_KEY_ID'] = aws_access_key_id
        os.environ['AWS_SECRET_ACCESS_KEY'] = aws_secret_access_key
        print(f"ACCESS KEY SECRET: {aws_secret_access_key}")
        sc = SparkContext()
        sc.setSystemProperty('com.amazonaws.services.s3.enableV4', 'true')
        hadoopConf = sc._jsc.hadoopConfiguration()
        hadoopConf.set("fs.s3a.awsAccessKeyId", aws_access_key_id)
        hadoopConf.set("fs.s3a.awsSecretAccessKey", aws_secret_access_key)
        hadoopConf.set("fs.s3a.endpoint", "s3.amazonaws.com")
        hadoopConf.set('fs.s3a.impl', 'org.apache.hadoop.fs.s3a.S3AFileSystem')
    
    else:
        sc = SparkContext()

    spark = SparkSession(sc)
    print(f'Reading from {source}!')
    print(f'Writing to {output_path}')
    run_etl(source, output_path, spark=spark)

def get_args():
    parser = argparse.ArgumentParser(description="Spark ETL CLI")
    parser.add_argument('--source', type=str, dest="source", help="source file type (should be name of module in pipeline/intake/sources/)", required=True)
    parser.add_argument('--output_path', type=str, dest="output_path", help="where to write output of etl", required=True)
    return parser.parse_args()


if __name__ == "__main__":
    args = get_args()
    spark_run_etl(args.source, args.output_path, os.environ.get('AWS_ACCESS_KEY_ID'), os.environ.get('AWS_SECRET_ACCESS_KEY'), local=os.environ.get('LOCAL', True))
