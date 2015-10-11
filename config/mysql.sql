/* *****************************************************************************
// Setup the preferences
// ****************************************************************************/
SET NAMES utf8 COLLATE 'utf8_unicode_ci';
SET foreign_key_checks = 1;
SET time_zone = '+00:00';
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
SET storage_engine = InnoDB;
SET CHARACTER SET utf8;

/* *****************************************************************************
// Remove old database
// ****************************************************************************/
DROP DATABASE IF EXISTS ninja;

/* *****************************************************************************
// Create new database
// ****************************************************************************/
CREATE DATABASE ninja DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;
USE ninja;

/* *****************************************************************************
// Create the tables
// ****************************************************************************/
CREATE TABLE user_status (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    PRIMARY KEY (id)
);

INSERT INTO `user_status` (`id`, `status`, `created_at`, `updated_at`, `deleted`) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(3, 'not verified', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(4, 'reverify', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);

CREATE TABLE user (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password CHAR(60) NOT NULL,
    
    status_id TINYINT(1) UNSIGNED NOT NULL DEFAULT 3,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    UNIQUE KEY (email),
    CONSTRAINT `f_user_status` FOREIGN KEY (`status_id`) REFERENCES `user_status` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);

CREATE TABLE user_verification (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    token VARCHAR(6) NOT NULL,
    
    user_id INT(10) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    CONSTRAINT `f_photo_verification_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
	UNIQUE KEY (user_id),
    PRIMARY KEY (id)
);

CREATE TABLE user_login (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
	remote_address VARCHAR(50) NOT NULL,
    user_id INT(10) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    CONSTRAINT `f_user_login_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);

CREATE TABLE email_verification (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    token VARCHAR(32) NOT NULL,
    
    user_id INT(10) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    CONSTRAINT `f_email_verification_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);

CREATE TABLE photo_status (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    PRIMARY KEY (id)
);

INSERT INTO `photo_status` (`id`, `status`, `created_at`, `updated_at`, `deleted`) VALUES
(1, 'verified',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'unverified', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(3, 'rejected', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(4, 'reverify', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);

CREATE TABLE photo (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    path VARCHAR(14) NOT NULL,
	
    user_id INT(10) UNSIGNED NOT NULL,
	note VARCHAR(255) NOT NULL DEFAULT "",
	
	initial TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
	
	status_id TINYINT(1) UNSIGNED NOT NULL DEFAULT 2,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    CONSTRAINT `f_photo_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `f_photo_status` FOREIGN KEY (`status_id`) REFERENCES `photo_status` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);

CREATE TABLE site (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    name VARCHAR(50) NOT NULL,
	url VARCHAR(255) NOT NULL,
	profile VARCHAR(255) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
	UNIQUE KEY (name),
    PRIMARY KEY (id)
);

INSERT INTO `site` (`id`, `name`, `url`, `profile`, `created_at`, `updated_at`, `deleted`) VALUES
(1, 'OKCupid',	'http://www.okcupid.com/',	'http://www.okcupid.com/profile/:name',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'ChristianMingle',	'http://www.christianmingle.com/',	'http://www.christianmingle.com/view-profile.html?u=:name', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);

CREATE TABLE username (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    name VARCHAR(100) NOT NULL,
	
    user_id INT(10) UNSIGNED NOT NULL,
	site_id INT(10) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    CONSTRAINT `f_username_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `f_username_site` FOREIGN KEY (`site_id`) REFERENCES `site` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
	UNIQUE KEY (site_id, user_id),
	UNIQUE KEY (name, site_id),
    PRIMARY KEY (id)
);

CREATE TABLE role_level (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    name VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    PRIMARY KEY (id)
);

INSERT INTO `role_level` (`id`, `name`, `created_at`, `updated_at`, `deleted`) VALUES
(1, 'Administrator',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'User', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);

CREATE TABLE role (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
	user_id INT(10) UNSIGNED NOT NULL,
	
    level_id TINYINT(1) UNSIGNED NOT NULL DEFAULT 2,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    UNIQUE KEY (user_id),
    CONSTRAINT `f_role_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT `f_role_level` FOREIGN KEY (`level_id`) REFERENCES `role_level` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);

CREATE TABLE ethnicity (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    user_id INT(10) UNSIGNED NOT NULL,
	
	type_id TINYINT(1) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
	
	CONSTRAINT `f_ethnicity_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);

CREATE TABLE ethnicity_type (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    name VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    PRIMARY KEY (id)
);

INSERT INTO `ethnicity_type` (`id`, `name`, `created_at`, `updated_at`, `deleted`) VALUES
(1, 'Asian',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'Middle Eastern', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(3, 'Black', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(4, 'Native American', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(5, 'Indian', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(6, 'Pacific Islander', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(7, 'Hispanic / Latin', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(8, 'White', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(9, 'Other', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);

CREATE TABLE demographic (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,

	user_id INT(10) UNSIGNED NOT NULL,
    
    birth_month TINYINT(2) UNSIGNED  NOT NULL,
	birth_day TINYINT(2) UNSIGNED  NOT NULL,
	birth_year SMALLINT(4) UNSIGNED  NOT NULL,
	
	gender CHAR(1) NOT NULL,
	
	height_feet TINYINT(2) UNSIGNED NOT NULL,
	height_inches TINYINT(2) UNSIGNED NOT NULL,
	
	weight SMALLINT(3) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    UNIQUE KEY (user_id),
    CONSTRAINT `f_demographic_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);

CREATE TABLE tracking_url (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	
	user_id INT(10) UNSIGNED NOT NULL,
    method VARCHAR(25) NOT NULL,
    url VARCHAR(2048) NOT NULL,
	remote_address VARCHAR(50) NOT NULL,
	referer VARCHAR(2048) NOT NULL,
	user_agent VARCHAR(8192) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id)
);

CREATE TABLE tracking_api (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	
	user_id INT(10) UNSIGNED NOT NULL,
    method VARCHAR(25) NOT NULL,
    url VARCHAR(2048) NOT NULL,
	remote_address VARCHAR(50) NOT NULL,
	referer VARCHAR(2048) NOT NULL,
	user_agent VARCHAR(8192) NOT NULL,
	lookup_user_id INT(10) UNSIGNED NOT NULL,
	verified TINYINT(1) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id)
);

CREATE TABLE api_authentication (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    userkey VARCHAR(32) NOT NULL,
	token VARCHAR(32) NOT NULL,
    
    user_id INT(10) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    CONSTRAINT `f_api_authentication_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);