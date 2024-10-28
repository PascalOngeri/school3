<?php 
// DB credentials.
define('DB_HOST','173.249.20.229');
define('DB_USER','remote');
define('DB_PASS','Qwerty245!');
define('DB_NAME','schoolsystem');
//define('DB_HOST','localhost');
//define('DB_USER','root');
//define('DB_PASS','@mesopotamia123');
//define('DB_NAME','schoolsystem');
//nnn

// Establish database connection.
try
{
$dbh = new PDO("mysql:host=".DB_HOST.";dbname=".DB_NAME,DB_USER, DB_PASS,array(PDO::MYSQL_ATTR_INIT_COMMAND => "SET NAMES 'utf8'"));
}
catch (PDOException $e)
{
exit("Error: " . $e->getMessage());
}



?><!--  Orginal Author Name: Mayuri K. 
 for any PHP, Codeignitor, Laravel OR Python work contact me at mayuri.infospace@gmail.com  
 Visit website : www.mayurik.com -->  
