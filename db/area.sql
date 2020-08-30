CREATE TABLE area (
 id  SERIAL PRIMARY KEY,
 zip VARCHAR(7) NOT NULL,
 pref VARCHAR(100) ,
 city VARCHAR(100) ,
 town VARCHAR(100) ,
 pref_kana VARCHAR(100) ,
 city_kana VARCHAR(100) ,
 town_kana VARCHAR(100)
)
