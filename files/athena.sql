-- create table
CREATE EXTERNAL TABLE IF NOT EXISTS records (
  ID BIGINT,
  Number INT,
  Type INT,
  UserID BIGINT,
  UserType INT,
  UserLogin STRING,
  CreatedAt TIMESTAMP
  )
  ROW FORMAT SERDE 'org.apache.hive.hcatalog.data.JsonSerDe'
  WITH SERDEPROPERTIES ( 'serialization.format' = '1' )
  LOCATION 's3://collector-data-bucket/';