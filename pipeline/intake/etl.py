import pkg_resources
from pyspark.sql import SparkSession
from pyspark import SparkFiles
from pyspark.sql import Row
import yaml

def filter_helper(partitionData, header, ignore_symbol): # TODO - move into file-specific module
    """
    Filter Lambda for RDD.
    Used to filter out all rows that do not match 
    expected row structure.

    :params partitionData - generator of strings. each element
        is a row of a source file:
            ex: '1974,5,19,1974.3795,333.37,5,-999.99,-999.99,50.40'
    :params header (string) - expected header line of file
    :params ignore_symbol (string) - indicator that a line should be filtered
        ex: '#'
    
    :returns generator of filtered source file lines.
    """
    for row in partitionData:
        if ignore_symbol != row[0] and row != header:
            yield row


def structure_as_row(partitionData, header_keys):
    """
    Map unstructured data (string) to structured representation
    with keys and values.

    :params partitionData - generator of strings.
    :params header_keys (list) - ordered list of headers in row

    :returns generator of Row records
    """
    for row in partitionData:
        values = row.split(',')
        if len(values) != len(header_keys):
            raise RuntimeError("Error! Number of Header Keys Does Not Match Number of Field Values!")
        record = {}
        for idx, key in enumerate(header_keys):
            record[key] = values[idx]
        
        yield Row(**record)


def run_etl(source, output_path, spark=None):
    """
    Run Spark ETL of source file.

    :params source (string) - name of source type (should be module in intake/sources/)
    :param output_path (string) - where to write parquet output
    :params spark - spark context
    """
    if not spark:
        spark = SparkSession.builder.getOrCreate()

    config = yaml.safe_load(pkg_resources.resource_stream(f'intake.sources.{source}', f'{source}_config.yml'))
    file_path = config['source']
    header_keys = config['header_keys']
    ignore_symbol = config['ignore_symbol']

    spark.sparkContext.addFile(file_path)
    data_path = SparkFiles.get(file_path.split('/')[-1])
    rdd = spark.sparkContext.textFile(data_path)

    # Use mapPartitions for structuring rows to only load
    # keys once per partition. Alternatively, we can consider
    # broadcasting the header_keys to workers...
    df = rdd.mapPartitions(lambda partition: filter_helper(partition, header=','.join(header_keys), ignore_symbol=ignore_symbol)) \
        .mapPartitions(lambda partition: structure_as_row(partition, header_keys)).toDF()
    df.write.parquet(output_path)