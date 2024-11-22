<?php 
// Include Composer's autoloader to load the dotenv library
require 'vendor/autoload.php';

// Specify the path to the .env file in the 'includes' directory
$dotenv = Dotenv\Dotenv::createImmutable(__DIR__ . '/includes');
$dotenv->load();

// Fetch database credentials from the .env file
$DB_HOST = $_ENV['DB_HOST'];
$DB_USER = $_ENV['DB_USER'];
$DB_PASS = $_ENV['DB_PASS'];
$DB_NAME = $_ENV['DB_NAME'];

// Establish database connection using PDO
try {
    // Using PDO to connect to the database with the credentials from the .env file
    $dbh = new PDO("mysql:host=$DB_HOST;dbname=$DB_NAME", $DB_USER, $DB_PASS, 
        array(PDO::MYSQL_ATTR_INIT_COMMAND => "SET NAMES 'utf8'"));
} catch (PDOException $e) {
    // If there's an error with the connection, display it
    exit("Error: " . $e->getMessage());
}
?>

<!-- Original Author Info -->
<!-- 
    Original Author Name: Mayuri K. 
    For any PHP, CodeIgniter, Laravel, or Python work contact me at mayuri.infospace@gmail.com  
    Visit website: www.mayurik.com
-->

