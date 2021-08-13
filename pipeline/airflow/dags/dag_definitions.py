from airflow import DAG

from dag_templates.intake_dag import create_intake_dag

globals()["co2_weekly_mlo"] =  create_intake_dag("co2_weekly_mlo")

globals()["ch4_mm_gl"] =  create_intake_dag("ch4_mm_gl")
