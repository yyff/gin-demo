

create table if not exists orders (
   id  INT(10) unsigned NOT NULL AUTO_INCREMENT,
   	product_id INT(10) unsigned NOT NULL,
   	user_id INT(10) unsigned NOT NULL,
   PRIMARY KEY (id)
)ENGINE=InnoDB  DEFAULT CHARSET=utf8;
