function sshLogin (ip, password, username) {
    const { exec } = require('child_process');
    
    exec(`sshpass -p ${password} ssh ${username}@${ip}`, (err, stdout, stderr) => {
        if (err) {
            console.error(err);
            return;
        }
        console.log(stdout);
    });
}

function sshCommand (command) {
    const { exec } = require('child_process');
    exec(command, (err, stdout, stderr) => {
        if (err) {
            console.error(err);
            return;
        }
        console.log(stdout);
    });
}

function runCrawler (credentials) {
    console.log(credentials);
    for (let i = 0; i < credentials.length; i++) {
        sshLogin(credentials[i].ip, credentials[i].password, credentials[i].username);
        console.log('Logged in');
        sshCommand(`cd ${credentials[i].path_to_crawler}`);
        console.log('Changed directory');
        if (credentials[i].isHost) {
            sshCommand(`cd ${credentials[i].path_to_task_manager}`);
            sshCommand(`go run manager.go`);
            sshCommand(`cd ${credentials[i].path_to_crawler}`);
            sshCommand(`go run main.go master ${credentials[i].ip}`);
            console.log('Host');
            console.log(credentials[i]);
        } else {
            sshCommand(`go run main.go worker ${credentials[i].ip}`);
            // console.log('Worker');
            // console.log(credentials[i]);
        }
    }
}

export default runCrawler;