drop database if exists pikopos;
create database pikopos;
use pikopos;

-- create "company, company_detail, outlet"
create table company (
  id int not null auto_increment,
  username varchar(16) not null,
  name varchar(32) not null,
  status enum('free', 'free-ads', 'paid-01') not null,
  primary key (id),
  unique (username)
);

create table company_detail (
  id int not null auto_increment,
  company_id int not null,
  address varchar(128),
  city varchar(64),
  province varchar(32),
  postal varchar(5),
  lat decimal(9,6),
  lon decimal(9,6),
  phone_number varchar(16),
  whatsapp varchar(16),
  email varchar(48),
  website varchar(64),
  twitter varchar(48),
  facebook varchar(48),
  primary key (id),
  foreign key (company_id) references company(id)
);

create table outlet (
  company_id int not null,
  id int not null auto_increment,
  name varchar(32),
  address varchar(128),
  city varchar(64),
  province varchar(32),
  postal varchar(5),
  lat decimal(9,6),
  lon decimal(9,6),
  phone_number varchar(16),
  whatsapp varchar(16),
  status enum('inactive', 'active') not null,
  primary key (id),
  foreign key (company_id) references company(id)
);

-- create "role, privilege"
create table role (
  company_id int not null,
  id int not null auto_increment,
  name varchar(16) not null,
  status enum('inactive', 'active') not null,
  primary key (id),
  foreign key (company_id) references company(id)
);

create table privilege (
  id int not null auto_increment,
  tag varchar(128) not null, -- max 128 char for privilege tag
  employee_required tinyint(1) not null,
  primary key (id),
  unique (tag)
);

create table role_privilege (
  company_id int not null,
  id int not null auto_increment,
  role_id int not null,
  privilege_id int not null,
  primary key (id),
  unique (role_id, privilege_id),
  foreign key (role_id) references role(id),
  foreign key (privilege_id) references privilege(id),
  foreign key (company_id) references company(id)
);

-- create "employee, employee_login, employee_outlet"
create table employee (
  company_id int not null,
  id int not null auto_increment,
  full_name varchar(32) not null,
  -- email varchar(48) not null,
  -- phone_number varchar(16) not null,
  -- password varchar(64) not null, -- hashed
  role_id int not null,
  pin varchar(64), -- hashed
  status enum('unverified', 'inactive', 'active') not null,
  primary key (id),
  foreign key (role_id) references role(id),
  foreign key (company_id) references company(id)
);

-- create table employee_register (
--   id int not null auto_increment,
--   employee_id int not null,
--   email varchar(48) not null,
--   phone_number varchar(16) not null,
--   otp_code varchar(64) not null,
--   expired_at datetime not null,
--   primary key (id),
--   foreign key (employee_id) references employee(id)
-- );

create table employee_login (
  id int not null auto_increment,
  employee_id int not null,
  provider enum('google') not null,
  id_from_provider varchar(255) not null,
  additional_data json,
  primary key (id),
  foreign key (employee_id) references employee(id)
);

create table employee_outlet (
  id int not null auto_increment,
  employee_id int not null,
  outlet_id int not null,
  primary key (id),
  unique (employee_id, outlet_id),
  foreign key (employee_id) references employee(id),
  foreign key (outlet_id) references outlet(id)
);

-- create "unit, unit_conversion"
create table unit (
  id int not null auto_increment,
  name varchar(24) not null,
  alias varchar(6) not null,
  primary key (id)
);

create table unit_conversion (
  id int not null auto_increment,
  unit_id_from int not null,
  unit_id_to int not null,
  multiplier decimal(6,2) not null,
  primary key (id),
  unique (unit_id_from, unit_id_to)
);

-- create "item_category, item_variant_group, item"
create table item_category (
  company_id int not null,
  id int not null auto_increment,
  name varchar(24) not null,
  status enum('inactive', 'active') not null,
  primary key (id),
  foreign key (company_id) references company(id)
);

create table item_variant_group (
  company_id int not null,
  item_category_id int not null,
  id int not null auto_increment,
  name varchar(24) not null,
  status enum('inactive', 'active') not null,
  primary key (id),
  foreign key (company_id) references company(id),
  foreign key (item_category_id) references item_category(id)
);

create table item (
  company_id int not null,
  item_variant_group_id int not null,
  id int not null auto_increment,
  item_type enum('direct_item', 'raw_item', 'finshed_item') not null,
  name varchar(24) not null,
  sku varchar(24),
  inventory_stock int,
  unit_id int,
  inventory_cost bigint,
  price bigint,
  status enum('inactive', 'active') not null,
  primary key (id),
  foreign key (company_id) references company(id),
  foreign key (item_variant_group_id) references item_variant_group(id),
  foreign key (unit_id) references unit(id)
);

-- create "recipe, recipe_item, modifier_item"
create table recipe (
  company_id int not null,
  id int not null auto_increment,
  finished_item_id int not null,
  primary key (id),
  unique (finished_item_id),
  foreign key (company_id) references company(id),
  foreign key (finished_item_id) references item(id)
);

create table recipe_item (
  id int not null auto_increment,
  recipe_id int not null,
  raw_item_id int not null,
  primary key (id),
  unique (recipe_id, raw_item_id),
  foreign key (recipe_id) references recipe(id),
  foreign key (raw_item_id) references item(id)
);

create table modifier_item (
  company_id int not null,
  id int not null auto_increment,
  direct_finished_item_id int not null,
  raw_item_w_modifier_id int not null,
  primary key (id),
  unique (direct_finished_item_id, raw_item_w_modifier_id),
  foreign key (company_id) references company(id),
  foreign key (direct_finished_item_id) references item(id),
  foreign key (raw_item_w_modifier_id) references item(id)
);

-- create "item_price_history, item_inventory_batch"
create table item_price_history (
  id int not null auto_increment,
  item_id int not null,
  price bigint not null,
  time_start datetime not null,
  time_end datetime,
  primary key (id),
  foreign key (item_id) references item(id)
);

create table item_inventory_batch (
  id int not null auto_increment,
  item_id int not null,
  stock int not null,
  cost bigint not null,
  batch_time datetime,
  primary key (id),
  foreign key (item_id) references item(id)
);

-- create "purchase_order, purchase_order_item"
create table purchase_order (
  outlet_id int not null,
  id int not null auto_increment,
  purchase_number int not null,
  order_amount bigint not null,
  order_open_time datetime not null,
  order_finish_time datetime,
  order_by_employee_id int not null,
  order_status enum('open', 'canceled', 'fulfilled') not null,
  primary key (id),
  unique (purchase_number),
  foreign key (outlet_id) references outlet(id),
  foreign key (order_by_employee_id) references employee(id)
);

create table purchase_order_item (
  purchase_order_id int not null,
  id int not null auto_increment,
  item_id int not null,
  item_amount int not null,
  item_cost bigint not null,
  primary key (id),
  foreign key (purchase_order_id) references purchase_order(id),
  foreign key (item_id) references item(id)
);

-- create "transaction_order, transaction_order_item, transaction_order_item_modifier, payment_order"
create table transaction_order (
  outlet_id int not null,
  id int not null auto_increment,
  invoice_number int not null,
  order_amount bigint not null,
  order_open_time datetime not null,
  order_finish_time datetime,
  order_by_employee_id int not null,
  order_status enum('open', 'hold', 'canceled', 'issued', 'fulfilled') not null,
  primary key (id),
  unique (invoice_number),
  foreign key (outlet_id) references outlet(id),
  foreign key (order_by_employee_id) references employee(id)
);

create table transaction_order_item (
  transaction_order_id int not null,
  id int not null auto_increment,
  item_id int not null,
  item_amount int not null,
  item_price bigint not null,
  primary key (id),
  foreign key (transaction_order_id) references transaction_order(id),
  foreign key (item_id) references item(id)
);

create table transaction_order_item_modifier (
  transaction_order_item_id int not null,
  id int not null auto_increment,
  item_id int not null,
  modifier_amount int not null,
  modifier_price bigint not null,
  primary key (id),
  foreign key (transaction_order_item_id) references transaction_order_item(id),
  foreign key (item_id) references item(id)
);

create table payment_order (
  outlet_id int not null,
  transaction_order_id int not null,
  id int not null auto_increment,
  payment_type enum('cash') not null,
  payment_amount bigint not null,
  payment_time datetime,
  payment_status enum('processing', 'failed', 'success') not null,
  primary key (id),
  foreign key (outlet_id) references outlet(id),
  foreign key (transaction_order_id) references transaction_order(id)
);