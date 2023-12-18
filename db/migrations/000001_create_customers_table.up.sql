CREATE TABLE IF NOT EXISTS customer (
	id bigserial PRIMARY KEY,
	card_number varchar(16) not null,
	first_name varchar(50),
	home_address_3 varchar(50),
	home_address_4 varchar(50),
	home_zip_code varchar(20),
	collector varchar(20),
	concat_customer varchar(20),
	location_file varchar(50)
);
