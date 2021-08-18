CREATE TABLE ch4_mm_gl (
  year  int NOT NULL,
  month int NOT NULL,
  day int DEFAULT 1,
  date_decimal  real NOT NULL,
  average real  NOT NULL,
  average_unc  real  NOT NULL,
  trend  real  NOT NULL,
  trend_unc real  NOT NULL,
  YYYYMMDD  date  NOT NULL
);
CREATE UNIQUE INDEX idx_yyyymmdd ON ch4_mm_gl(YYYYMMDD);