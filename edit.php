
<?php
session_start();
include('dbconnection.php');
if (isset($_POST['submit'])) {
    // Process form submission
    // Update database with form data
        $fname = $_POST['fname'];
        $mname = $_POST['mname'];
        $lname = $_POST['lname'];
        $gender = $_POST['gender'];
        $dob = $_POST['dob'];
        $adm = $_POST['stuid'];
        $email = $_POST['stuemail'];
        $class = $_POST['class'];
        $phone = $_POST['connum'];
        $phone1 = $_POST['altconnum'];
        $address = $_POST['address'];
        $username = $_POST['uname'];
        $faname = $_POST['faname'];
        $maname = $_POST['maname'];
        $password = $_POST['password'];
     //   $image = $_FILES["image"]["name"];
 

    $query = mysqli_query($con, "UPDATE registration SET adm='$adm', fname='$fname', mname='$mname', lname='$lname', gender='$gender', class='$class', phone='$phone', phone1='$phone1', email='$email', address='$address',dob='$dob',username='$username'  ,password='$password'   ,faname='$faname',maname='$maname' WHERE adm='$adm'");
    if ($query) {
        echo '<script>alert("Update successful")</script>';
        echo "<script>window.location.href='manage-students.php'</script>";
        exit(); // Add exit() after redirect
    } else {
        echo '<script>alert("Something went wrong. Please try again")</script>';
        echo "<script>window.location.href='manage-students.php'</script>";
        exit(); // Add exit() after redirect
    }
}
?>
