CREATE TABLE `query_log` (
      `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
      `bindport` smallint(5) unsigned NOT NULL,
      `client` char(15) NOT NULL DEFAULT '',
      `client_port` smallint(5) unsigned NOT NULL,
      `server` char(15) NOT NULL DEFAULT '',
      `server_port` smallint(5) unsigned NOT NULL,
      `sql_type` varchar(30) NOT NULL DEFAULT 'Query',
      `sql_string` text,
      `create_time` datetime NOT NULL,
      PRIMARY KEY (`id`),
      KEY `idx_client` (`client`),
      KEY `idx_server` (`server`),
      KEY `idx_cretime` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=9945 DEFAULT CHARSET=utf8
