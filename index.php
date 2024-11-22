<?php
session_start();

error_reporting(0);

// Database connection
include('includes/dbconnection.php');
include('dbconnection.php');

// Set the desired port
$port = getenv('PORT') ?: 8000; // Default to 8000 if not specified

// Log the port being used
mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('system', 'Server running on port: $port')");

// Select data from the database
$selectQuery = "SELECT name, icon, iname FROM api";
$result = mysqli_query($con, $selectQuery);
if ($result && mysqli_num_rows($result) > 0) {
    $row = mysqli_fetch_assoc($result);
    $_SESSION['name'] = $row['name'];
    $_SESSION['icon'] = $row['icon'];
    $_SESSION['iname'] = $row['iname'];
    mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('system', 'Fetched data from api table')");
} else {
    mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('system', 'No data found in api table')");
}

function sendSMS($phoneNumber, $message) {
    include('includes/dbconnection.php');
    include('dbconnection.php');

    $url = 'https://sms-service.smsafrica.tech/message/send/transactional';

    // Retrieve the latest API key from the database
    $selectQuery = "SELECT apikey FROM api ORDER BY id DESC LIMIT 1";
    $result = mysqli_query($con, $selectQuery);
    if ($result && mysqli_num_rows($result) > 0) {
        $row = mysqli_fetch_assoc($result);
        $token = $row['apikey'];
    } else {
        echo "Error fetching latest API key from the database.\n";
        mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('system', 'Error fetching latest API key from the database')");
        exit();
    }

    $postData = json_encode([
        "message" => $message,
        "msisdn" => $phoneNumber,
        "sender_id" => "SMSAFRICA",
        "callback_url" => "https://callback.io/123/dlr"
    ]);

    $httpRequest = curl_init($url);
    curl_setopt($httpRequest, CURLOPT_POST, true);
    curl_setopt($httpRequest, CURLOPT_POSTFIELDS, $postData);
    curl_setopt($httpRequest, CURLOPT_TIMEOUT, 60);
    curl_setopt($httpRequest, CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($httpRequest, CURLOPT_HTTPHEADER, [
        "Content-Type: application/json",
        "api-key: $token"
    ]);

    $response = curl_exec($httpRequest);
    $httpCode = curl_getinfo($httpRequest, CURLINFO_HTTP_CODE);
    curl_close($httpRequest);

    if ($httpCode == 200) {
        mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('system', 'Sent SMS to $phoneNumber')");
    } else {
        mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('system', 'Failed to send SMS to $phoneNumber')");
    }
}

function checkAndSendFridaySMS($dbh, $con) {
    if (date('N') == 5) { // Check if today is Friday
        $today = date('Y-m-d');
        $check_sql = "SELECT * FROM logs WHERE user='system' AND activities='Sent Friday SMS to all users' AND date(date) = :today";
        $check_query = $dbh->prepare($check_sql);
        $check_query->bindParam(':today', $today, PDO::PARAM_STR);
        $check_query->execute();

        if ($check_query->rowCount() == 0) {
            $sql_phone_numbers = "SELECT MobileNumber FROM tbladmin";
            $query_phone_numbers = $dbh->prepare($sql_phone_numbers);
            $query_phone_numbers->execute();
            $results_phone_numbers = $query_phone_numbers->fetchAll(PDO::FETCH_OBJ);

            if ($query_phone_numbers->rowCount() > 0) {
                $message = "Friday message for all users."; // Replace with your message

                foreach ($results_phone_numbers as $result_phone) {
                    sendSMS($result_phone->MobileNumber, $message);
                }

                $insert_log_sql = "INSERT INTO logs(user, activities, date) VALUES('system', 'Sent Friday SMS to all users', :today)";
                $insert_log_query = $dbh->prepare($insert_log_sql);
                $insert_log_query->bindParam(':today', $today, PDO::PARAM_STR);
                $insert_log_query->execute();
            }
        } else {
            mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('system', 'Friday SMS already sent today')");
        }
    }
}

if (isset($_POST['login'])) {
    $username = $_POST['username'];
    $password = $_POST['password'];

    // Log the login attempt
    mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('$username', 'Attempted login')");

    // Check in tbladmin table
    $sql_admin = "SELECT ID, UserName FROM tbladmin WHERE UserName=:username AND Password=:password";
    $query_admin = $dbh->prepare($sql_admin);
    $query_admin->bindParam(':username', $username, PDO::PARAM_STR);
    $query_admin->bindParam(':password', $password, PDO::PARAM_STR);
    $query_admin->execute();
    $results_admin = $query_admin->fetchAll(PDO::FETCH_OBJ);

    if ($query_admin->rowCount() > 0) {
        foreach ($results_admin as $result_admin) {
            $_SESSION['sturecmsaid'] = $result_admin->ID;
            $_SESSION['username'] = $result_admin->UserName;
            mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('$username', 'Logged in successfully')");
        }

        // Set cookies if "Remember me" is checked
        if (!empty($_POST["remember"])) {
            setcookie("user_login", $_POST["username"], time() + (10 * 365 * 24 * 60 * 60));
            setcookie("userpassword", $_POST["password"], time() + (10 * 365 * 24 * 60 * 60));
            mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('$username', 'Checked remember me option')");
        } else {
            if (isset($_COOKIE["user_login"])) {
                setcookie("user_login", "", time() - 3600);
            }
            if (isset($_COOKIE["userpassword"])) {
                setcookie("userpassword", "", time() - 3600);
            }
        }
        $_SESSION['login'] = $_POST['username'];

        // Check and send Friday SMS
        checkAndSendFridaySMS($dbh, $con);

        echo "<script type='text/javascript'> document.location ='dashboard.php'; </script>";
    } else {
        // Check in tblparent table if not found in tbladmin
        $sql_parent = "SELECT id, adm, username, phone, password FROM registration WHERE username=:username AND password=:password";
        $query_parent = $dbh->prepare($sql_parent);
        $query_parent->bindParam(':username', $username, PDO::PARAM_STR);
        $query_parent->bindParam(':password', $password, PDO::PARAM_STR);
        $query_parent->execute();
        $results_parent = $query_parent->fetchAll(PDO::FETCH_OBJ);

        if ($query_parent->rowCount() > 0) {
            foreach ($results_parent as $result) {
                $_SESSION['sturecmsaid'] = $result->id;
                $_SESSION['adm'] = $result->adm;
                $_SESSION['username'] = $result->username;
                $_SESSION['phone'] = $result->phone;
                $_SESSION['password'] = $result->password;
            }

            if (!empty($_POST["remember"])) {
                setcookie("user_login", $_POST["username"], time() + (10 * 365 * 24 * 60 * 60));
                setcookie("userpassword", $_POST["password"], time() + (10 * 365 * 24 * 60 * 60));
            } else {
                if (isset($_COOKIE["user_login"])) {
                    setcookie("user_login", "", time() - 3600);
                }
                if (isset($_COOKIE["userpassword"])) {
                    setcookie("userpassword", "", time() - 3600);
                }
            }
            $_SESSION['login'] = $_POST['username'];

            // Check and send Friday SMS
            checkAndSendFridaySMS($dbh, $con);

            echo "<script type='text/javascript'> document.location ='pd.php'; </script>";
        } else {
            mysqli_query($con, "INSERT INTO logs(user, activities) VALUES('$username', 'Invalid login attempt')");
            echo "<script>alert('Invalid Details');</script>";
        }
    }
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <title>infinityschools</title>
    <!-- plugins:css -->
    <link rel="stylesheet" href="assets/vendors/simple-line-icons/css/simple-line-icons.css">
    <link rel="stylesheet" href="assets/vendors/flag-icon-css/css/flag-icon.min.css">
    <link rel="stylesheet" href="assets/vendors/css/vendor.bundle.base.css">
    <!-- endinject -->
    <!-- Plugin css for this page -->
    <!-- End plugin css for this page -->
    <!-- inject:css -->
    <!-- endinject -->
    <!-- Layout styles -->
    <link rel="stylesheet" href="assets/css/style.css">
    <style>
        .content-wrapper {
            background-image: url('assets/images/background.jpg');
            background-size: cover;
        }
    </style>
</head>
<body>
    <div class="container-scroller">                           
 <center><i style="color: green">infinity system</i></center>                      
</div>
<div class="container-scroller">
    <div class="container-fluid page-body-wrapper full-page-wrapper">
        <div class="content-wrapper d-flex align-items-center auth">
            <div class="row flex-grow">
                <div class="col-lg-4 mx-auto">
                    <div class="auth-form-light text-center p-5">
                        <div class="brand-logo">
                            <img src="<?php echo $_SESSION['icon'];?>">
                        </div>
                        <h4>Hello! let's get started</h4>
                        <h6 class="font-weight-light">Sign in to continue.</h6>
                        <form class="pt-3" id="login" method="post" name="login">
                            <div class="form-group">
                                <input type="text" class="form-control form-control-lg" placeholder="Enter your username" required="true" name="username" value="<?php if (isset($_COOKIE["user_login"])) { echo $_COOKIE["user_login"]; } ?>">
                            </div>
                            <div class="form-group">
                                <input type="password" class="form-control form-control-lg" placeholder="Enter your password" name="password" required="true" value="<?php if (isset($_COOKIE["userpassword"])) { echo $_COOKIE["userpassword"]; } ?>">
                            </div>
                            <div class="mt-3">
                                <button class="btn btn-success btn-block loginbtn" name="login" type="submit">Login</button>
                            </div>
                            <div class="my-2 d-flex justify-content-between align-items-center">
                                <div class="form-check">
                                    <label class="form-check-label text-muted">
                                        <input type="checkbox" id="remember" class="form-check-input" name="remember" <?php if (isset($_COOKIE["user_login"])) { ?> checked <?php } ?> /> Keep me signed in </label>
                                </div>
                                <a href="forgot.php" class="auth-link text-black">Forgot password?</a>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
        <!-- content-wrapper ends -->
    </div>
    <!-- page-body-wrapper ends -->
</div>
<!-- container-scroller -->
<!-- plugins:js -->
<script src="assets/vendors/js/vendor.bundle.base.js"></script>
<!-- endinject -->
<!-- Plugin js for this page -->
<!-- End plugin js for this page -->
<!-- inject:js -->
<script src="assets/js/off-canvas.js"></script>
<script src="assets/js/bootstrap.min.js"></script>
<!-- endinject -->
</body>
</html>
