# Project : <mark>Mini Google</mark>
Authors (team): <br>
<mark>Sviatoslav Stehnii https://github.com/sviatkohuet <br>
Taras Lysun https://github.com/taraslysun <br>
Dmytro Khamula https://github.com/hamuladm <br>
Bohdan Ozarko https://github.com/Compi-Craft
</mark><br>

## Prerequisites

<mark>Node.js and npm, Vite, Golang, ElasticSearch</mark>

### Usage


<mark>To run application you need to run client side and server side.<br>
First of all you have to execute server side, you can do this in app/api folder by running ```go run main.go```.<br>
Go to client folder in app/, run ```npm i``` to install all dependencies, then ```npm run dev``` to run frontend.<br>
Next step you have to fill your elasticsearch CloudId and APIKey in app/api/utils/vars.go.<br>
If needed you can start crawling websites using our crawler, make the previous step in crawler/utils/vars.go and run task_manager and crawler (before crawler running POST some links to task_manager) with command ```go run .```.<br>
Now you can search!</mark> 

### Example
<mark>![image](https://github.com/taraslysun/GOofySearch/assets/81622077/46d8f70c-5c6f-4cb3-9a8f-4b1918f17c01)

</mark>
