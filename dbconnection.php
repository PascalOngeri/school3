<?php
// Include Composer's autoloader to load the dotenv library
require 'vendor/autoload.php';

// Specify the path to the .env file in the 'includes' directory
$dotenv = Dotenv\Dotenv::createImmutable(__DIR__ . '/includes');
$dotenv->load();

// Fetch database credentials from the .env file
$dbHost = $_ENV['DB_HOST'];
$dbUser = $_ENV['DB_USER'];
$dbPass = $_ENV['DB_PASS'];
$dbName = $_ENV['DB_NAME'];

// Establish a connection to the MySQL database using the values from the .env file
$con = mysqli_connect($dbHost, $dbUser, $dbPass, $dbName);

// Check if the connection was successful
if (mysqli_connect_errno()) {
    // If there is an error, display it
    echo "Connection Failed: " . mysqli_connect_error();
    exit(); // Stop the script execution if the connection fails
} else {
    // Optionally, you can echo a success message (useful for debugging)
    echo "Successfully connected to the database!";
}
?>

