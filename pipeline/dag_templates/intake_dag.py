from airflow import DAG
from airflow.operators.python_operator import PythonOperator
from airflow.utils.dates import days_ago

from intake.spark_etl import spark_run_etl


def create_intake_dag(src_name, output_path):

    default_args = {
        'owner': 'tkobil',
        'start_date': days_ago(1),
        'depends_on_past': False,
        'email': ['testemail@gmail.com'], # sample email for github
        'email_on_failure': False
    }

    with DAG(
        'etl_{}'.format(src_name),
        default_args=default_args,
        description='ETL from NOAA for {}'.format(src_name),
        schedule_interval="@daily",
    ) as dag:

        get_data_task = PythonOperator(
            task_id="get_data_{}".format(src_name),
            python_callable=spark_run_etl,
            op_args=[src_name, output_path]
        )

        #TODO - add load to database task 
        # as opional task based on env variables
    
    return dag