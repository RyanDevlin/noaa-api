CREATE TABLE co2_weekly_mlo (
  year  int NOT NULL,
  month int NOT NULL,
  day int NOT NULL,
  date_decimal  real NOT NULL,
  average real  NOT NULL,
  ndays int NOT NULL,
  one_year_ago  real  NOT NULL,
  ten_years_ago  real  NOT NULL,
  increase_since_1800 real  NOT NULL,
  YYYYMMDD  date  NOT NULL
);
CREATE UNIQUE INDEX idx_yyyymmdd ON co2_weekly_mlo(YYYYMMDD);