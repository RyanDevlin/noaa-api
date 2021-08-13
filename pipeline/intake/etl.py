import pkg_resources
from pyspark.sql import SparkSession
from pyspark import SparkFiles
from pyspark.sql import Row
import yaml
from datetime import datetime

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
        if ignore_symbol != row[0] and row.replace(' ', '_') != header:
            yield row


def structure_as_row(partitionData, header_keys):
    """
    Map unstructured data (string) to structured representation
    with keys and values.

    :params partitionData - generator of strings.
    :params header_keys (dict) - ordered dict of headers in row, with correspondign
                                data type expressions as their values

    :returns generator of Row records
    """
    for row in partitionData:
        values = row.split(',')
        if len(values) != len(header_keys):
            raise RuntimeError("Error! Number of Header Keys Does Not Match Number of Field Values!")

        record = {}
        for idx, key in enumerate(list(header_keys.keys())):
            data_type_expr = header_keys[key]
            try:
                record[key] = eval(data_type_expr)(values[idx])
            except Exception as ex:
                print(values)
                print(idx)
                raise ex
        
        yield Row(**record)


def create_yyyymmdd_index(row_dict):
    """
    Create a YYYYMMDD Key in row_dict,
    as return as a spark sql Row.
    YYYYMMDD is meant to be a value that can 
    be indexed on in a database...

    :params row_dict (dict)
    :returns pyspark.sql.Row
    """
    year = str(row_dict['year'])
    month = str(row_dict.get('month', '01')).zfill(2)
    day = str(row_dict.get('day', '01')).zfill(2)
    row_dict['YYYYMMDD'] = datetime.strptime(f'{year}{month}{day}', '%Y%m%d')
    return Row(**row_dict)


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
    # TODO - refactor column renames/yyyymmdd index creation as add more data sources...
    df = rdd.mapPartitions(lambda partition: filter_helper(partition, header=','.join(list(header_keys.keys())), ignore_symbol=ignore_symbol)) \
        .mapPartitions(lambda partition: structure_as_row(partition, header_keys)) \
        .map(lambda Row: create_yyyymmdd_index(Row.asDict())).toDF() \
        .withColumnRenamed("1_year_ago", "one_year_ago") \
        .withColumnRenamed("10_years_ago", "ten_years_ago") \
        .withColumnRenamed("decimal", "date_decimal")
    df.write.mode("overwrite").parquet(output_path) # Always overwrite with latest dataset