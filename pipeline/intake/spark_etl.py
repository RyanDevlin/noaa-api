import argparse
from pyspark.sql import SparkSession

from intake.etl import run_etl


def spark_run_etl(source, output_path):
    spark = SparkSession.builder.getOrCreate()
    run_etl(source, output_path, spark=spark)

def get_args():
    parser = argparse.ArgumentParser(description="Spark ETL CLI")
    parser.add_argument('--source', help="source file type (should be name of module in pipeline/intake/sources/)", required=True)
    parser.add_argument('--output_path', help="where to write output of etl", required=True)
    return parser.parse_args()


if __name__ == "__main__":
    args = get_args()
    spark_run_etl(args.source, args.output_path)
