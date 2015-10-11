# Verified.ninja Web Application

This project is on an indefinite hiatus. The website, https://verified.ninja, is no longer available online
so all the code is now open source.

The Chrome browser extension can be found [here](https://github.com/verifiedninja/chromeext).

## Overview

One of the predominate issues in online dating today is people don't look like their photos.
It's very easy for one person to impersonate another, upload old photos, and even use camera tricks/Photoshop
to make themselves look different from how they look today.

We designed a service where members of a dating website can create an account and then upload photos of themselves
that meet strict requirements. The member must follow these steps:

1. Create an account with first name, last name, email, and password. The email must be verified every 30 days so it
must be real.
2. Verify email address by clicking on the URL in the email.
3. Enter in demographic information with birthday, gender, height, weight, and ethnicity.
4. Upload a private photo which is used only for verification purposes by the administrators.
5. Enter dating website usernames to ensure no other person can claim them.
6. Upload public photos for verification.
7. Upload verified public photos with Verified.ninja logo to the dating website to show they are verified.

The web application is coded in Golang. There is also publicly available code for a Chrome extension.

The web application was created from this template: https://github.com/josephspurrier/gowebapp

## Hiatus

The service is pretty much complete, but there are still a few things that need 
to be added like a 90 day automatic private photo unverification which is pretty
easy. The route from the cron job should be protected so it's not public. The
brute force protection needs to be moved outside the application. There are a few
more things, but for the most part, it was a fully working service.

## Quick Start Testing

The application lived on Amazon Web Services, but it can be tested locally.

To download, run the following command:

~~~
go get github.com/verifiedninja/webapp
~~~

Start MySQL and import config/mysql.sql to create the database and tables.

Open config/config.json and edit the Database section so the connection information matches your MySQL instance.

Build and run from the root directory. Open your web browser to: http://localhost. You should see the home page.

Navigate to the register page. Create a new user and 
then go to the URL that displays in the console to verify the email address or 
set the user.status_id field to 1 through MySQL. You will then able to login.
You can also change the role.level_id field to 1 to set the user as an Administrator 
with verification capabilities.

## Server Configuration

Add this line to cron to require email verification after 30 days:

~~~
0 8 * * * wget -O /dev/null -q https://verified.ninja/cron/notifyemailexpire
~~~

## Verification Process

These are the requirements for the private photos that are showed to the member. 
It's the responsibility of administrator to verify the photo against these rules.

These are the requirements for the private photos:
* Hold a piece of white paper with the generated characters in a dark color so they can be read easily
* No selfies - the photo must be taken by another person
* Clearly shows your face, torso, arms, legs, and feet
* Smile - it's not a mug shot so don't look like it
* Taken in a well-lit area
* No mirror or bathroom photos
* No baggy clothing
* No other people in the photo
* No post processing - please don't apply any filters, affix any borders, or change the colors
* Photo type is JPG/JPEG, PNG, or PNG
* Must be atleast 300x300 pixels
* Must be less than 5MB (megabytes)
* No nudity or explicit content

The administrator can reject the photo with an explanation. Once a public photo 
is verified, the member can upload public photos for verification.

These are the requirements for the public photos:
* Taken in a well-lit area
* Clearly shows your face
* No mirror or bathroom photos
* No other people in the photo
* No post processing - please don't apply any filters, affix any borders, or change the colors
* Photo type is JPG/JPEG, PNG, or PNG
* Must be atleast 300x300 pixels
* Must be less than 5MB (megabytes)
* No nudity or explicit content

It's the responsibility of the administrator to verify the public photo meets the rules
AND looks like the private photo. Once the administrator approved the photo, the system 
adds a green ninja logo to the top right corner of the photo. The final step required 
to reach Verified.ninja status is to add the username of your dating website to 
the profile so no other person can claim the username on our website.

Once the member has the status of a Verified.ninja, he or she must maintain the 
status by reconfirming their email address every 30 days and uploading a new 
private photo every 90 days. If any of the public photos no long match the new 
private photo, they should be rejected by the Administrator.

When a member reaches the status of Verified.ninja, he or she can upload the 
photos with the green ninja logo to the dating website and then add a link to 
their public profile with text that says they are a Verified.ninja:

## Checking a Member's Verification Status

There are two ways for other members to confirm whether a person has the status 
of a Verified.ninja:

1. The first method is through the publicly available "Verify a User" page. 
The member chooses the dating website from the dropdown list and then types in the 
username of the member. Once the member hits the Verify button, a new page displays 
with the member's status and an explanation about our service.
2. The second method is much easier and quicker. We've developed a [Chrome browser 
extension](https://github.com/verifiedninja/chromeext) that will scan the dating website and insert a green ninja logo next to 
each verified username or a red ninja logo for those that have not verified with 
our service.

The two methods above protect against stolen photos and outdated accounts. One 
of the benefits of the service is members with the status of Verified.ninja usually 
get more responses on dating websites because his or her photos are verified.

We currently support the following dating websites: OKCupid and ChristianMingle.

## Screenshots

Public Home:

![Image of Public Home](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/home.PNG)

Register:

![Image of Register](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/register.PNG)

Login:

![Image of Register](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/login.PNG)

Demographic:

![Image of Demographic](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/demographic.PNG)

Private Photo Upload:

![Image of Private Upload](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/private.PNG)

Waiting for Private Photo Verification:

![Image of Private Photo Verification](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/waiting.PNG)

Verified Private Photo:

![Image of Private Photo](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/verifiedprivate.PNG)

Dating Usernames:

![Image of Dating Usernames](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/usernames.PNG)

Public Photo Upload:

![Image of Public Upload](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/public.PNG)

Verified Public Photo:

![Image of Public Upload](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/privateprofile.PNG)

Waiting for Public Photo Verification:

![Image of Public Photo Verification](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/privateprofile.PNG)

Public Profile

![Image of Public Profile](https://raw.githubusercontent.com/verifiedninja/webapp/master/github_screenshots/verifiedpublic.PNG)

## Feedback

All feedback is welcome. This code is no longer maintained, but can be used as a real-world example of how to build a web application
using Golang.