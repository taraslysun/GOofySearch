// import React, { useState } from "react";
// import axios from "axios";

// function CrawlerRunner() {
//     const [credentials, setCredentials] = useState([]);
//     const [formData, setFormData] = useState({
//         ip: '',
//         username: '',
//         password: '',
//         path_to_crawler: '',
//         is_host: false,
//         path_to_task_manager: '',
//     });

//     const [esCredentials, setEsCredentials] = useState({
//         cloud_id: '',
//         api_key: '',
//     });

//     const handleChange = (e) => {
//         const { name, value, type, checked } = e.target;
//         setFormData({
//             ...formData,
//             [name]: type === 'checkbox' ? checked : value
//         });
//     };

//     const handleSubmit = () => {
//         setCredentials(credentials.concat([formData]));
//         setFormData({
//             ip: '',
//             username: '',
//             password: '',
//             path_to_crawler: '',
//             path_to_task_manager: '',
//             is_host: false,
//         });
//     };

//     const handleDelete = (index) => {
//         const updatedCredentials = credentials.filter((_, i) => i !== index);
//         setCredentials(updatedCredentials);
//     };

//     const runCrawlers = async () => {
//         let host = credentials.find(credential => credential.is_host);
//         let workers = credentials.filter(credential => !credential.is_host);
        
//         let data = {
//             ip: host.ip,
//             username: host.username,
//             password: host.password,
//             path_to_crawler: host.path_to_crawler,
//             path_to_task_manager: host.path_to_task_manager,
//             is_host: true,
//             id: 1,
//             host_ip: host.ip,
//             worker_num: workers.length,
//         }
//         // console.log("Host Data: ", data);
//         await axios.post("http://localhost:3000/api/execute_ssh", data);

//         for (let i = 2; i < workers.length + 2; i++) {
//             data = {
//                 ip: workers[i - 2].ip,
//                 username: workers[i - 2].username,
//                 password: workers[i - 2].password,
//                 path_to_crawler: workers[i - 2].path_to_crawler,
//                 path_to_task_manager: workers[i - 2].path_to_task_manager,
//                 is_host: false,
//                 id: i,
//                 host_ip: host.ip,
//                 worker_num: workers.length,
//             };
//             // console.log("Worker Data: ", data);
//             await axios.post("http://localhost:3000/api/execute_ssh", data);
//         }



//     };

//     return (
//         <div style={{
//             display: 'flex',
//             flexDirection: 'column',
//             justifyContent: 'center',
//             alignItems: 'center',
//             height: '100vh',
//             width: '100vw',
//             backgroundColor: 'gray',
//             color: 'black',
        
//         }}>
//             <div style={{
//                 display: 'flex',
//                 flexDirection: 'column',
//                 justifyContent: 'center',
//                 alignItems: 'center',
//                 height: '30%',
//                 width: '50%',
//                 margin: 'auto',
//                 backgroundColor: 'lightgray',
//                 color: 'black',
//             }}>
//                 <input type="text" name="ip" placeholder="IP Address" value={formData.ip} onChange={handleChange} />
//                 <input type="text" name="username" placeholder="Username" value={formData.username} onChange={handleChange} />
//                 <input type="password" name="password" placeholder="Password" value={formData.password} onChange={handleChange} />
//                 <input type="text" name="path_to_crawler" placeholder="Path to Crawler" value={formData.path_to_crawler} onChange={handleChange} />
//                 <label>
//                     <input type="checkbox" name="is_host" checked={formData.is_host} onChange={handleChange} />
//                     Is Host
//                 </label>
//                 <input type="text" name="path_to_task_manager" placeholder="Path to Task Manager" value={formData.path_to_task_manager} onChange={handleChange} />
//                 <button onClick={handleSubmit}>Add Credentials</button>
//             </div>

//             <div style={{
//                 display: 'flex',
//                 flexDirection: 'row',
                
//                 justifyContent: 'center',
//                 alignItems: 'center',
//                 height: '30vh',
//                 width: '50vw',
//                 margin: 'auto',
//                 backgroundColor: 'lightgray',
//                 color: 'black',
//             }}>
//                 {credentials.map((credential, index) => (
//                     <div key={index} style={{
//                         display: 'flex',
//                         flexDirection: 'column',
//                         justifyContent: 'center',
//                         alignItems: 'center',
//                         height: '30vh',
//                         width: '30%',
//                         margin: 'auto',
//                         backgroundColor: 'lightgray',
//                         color: 'black',
//                     }}>
//                         <ul>
//                             <li>IP: {credential.ip}</li>
//                             <li>Username: {credential.username}</li>
//                             <li>Path to Crawler: {credential.path_to_crawler}</li>
//                             <li>Is Host: {credential.is_host.toString()}</li>
//                             <li>Path to Task Manager: {credential.path_to_task_manager}</li>
//                         </ul>
//                         <button onClick={() => handleDelete(index)}>Delete</button>

//                     </div>
//                 ))}
//             </div>

//             <div style={{
//                 display: 'flex',
//                 flexDirection: 'column',
//                 justifyContent: 'center',
//                 alignItems: 'center',
//                 height: '30vh',
//                 width: '50vw',
//                 margin: 'auto',
//                 backgroundColor: 'lightgray',
//                 color: 'black',
//             }}>
//                 <button onClick={() => 
//                     runCrawlers(credentials)
//                 } style={{
//                     width: '100%',
//                     height: '100%',
//                     margin: 'auto',
//                     backgroundColor: 'lightblue',
//                     color: 'black',
                
//                 }}>Run Crawler System</button>
//             </div>
//         </div>
//     );
// }

// export default CrawlerRunner;



import React, { useState } from "react";
import axios from "axios";

function CrawlerRunner() {
    const [credentials, setCredentials] = useState([]);
    const [formData, setFormData] = useState({
        ip: '',
        username: '',
        password: '',
        path_to_crawler: '',
        is_host: false,
        path_to_task_manager: '',
    });
    const [esCredentials, setEsCredentials] = useState({
        cloud_id: '',
        api_key: '',
    });

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        if (name === "cloud_id" || name === "api_key") {
            setEsCredentials({
                ...esCredentials,
                [name]: value,
            });
        } else {
            setFormData({
                ...formData,
                [name]: type === 'checkbox' ? checked : value
            });
        }
    };

    const handleSubmit = () => {
        setCredentials(credentials.concat([formData]));
        setFormData({
            ip: '',
            username: '',
            password: '',
            path_to_crawler: '',
            path_to_task_manager: '',
            is_host: false,
        });
    };

    const handleEsSubmit = async () => {
        await axios.post("http://localhost:3000/api/api_credentials", esCredentials);
        alert("Elasticsearch credentials updated!");
    };

    const handleDelete = (index) => {
        const updatedCredentials = credentials.filter((_, i) => i !== index);
        setCredentials(updatedCredentials);
    };

    const runCrawlers = async () => {
        let host = credentials.find(credential => credential.is_host);
        let workers = credentials.filter(credential => !credential.is_host);
        
        let data = {
            ip: host.ip,
            username: host.username,
            password: host.password,
            path_to_crawler: host.path_to_crawler,
            path_to_task_manager: host.path_to_task_manager,
            is_host: true,
            id: 1,
            host_ip: host.ip,
            worker_num: workers.length,
        }
        await axios.post("http://localhost:3000/api/execute_ssh", data);

        for (let i = 2; i < workers.length + 2; i++) {
            data = {
                ip: workers[i - 2].ip,
                username: workers[i - 2].username,
                password: workers[i - 2].password,
                path_to_crawler: workers[i - 2].path_to_crawler,
                path_to_task_manager: workers[i - 2].path_to_task_manager,
                is_host: false,
                id: i,
                host_ip: host.ip,
                worker_num: workers.length,
            };
            await axios.post("http://localhost:3000/api/execute_ssh", data);
        }
    };

    return (
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
            width: '100vw',
            backgroundColor: 'gray',
            color: 'black',
        }}>
            {/* Elasticsearch Credentials Form */}
            <div style={{
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: 'center',
                height: '10%',
                width: '50%',
                margin: '10px auto',
                backgroundColor: 'lightgray',
                color: 'black',
            }}>
                <input type="text" name="cloud_id" placeholder="Elasticsearch Cloud ID" value={esCredentials.cloud_id} onChange={handleChange} />
                <input type="text" name="api_key" placeholder="Elasticsearch API Key" value={esCredentials.api_key} onChange={handleChange} />
                <button onClick={handleEsSubmit}>Set Elasticsearch Credentials</button>
            </div>
//             <div style={{
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: 'center',
                height: '30%',
                width: '50%',
                margin: 'auto',
                backgroundColor: 'lightgray',
                color: 'black',
            }}>
                <input type="text" name="ip" placeholder="IP Address" value={formData.ip} onChange={handleChange} />
                <input type="text" name="username" placeholder="Username" value={formData.username} onChange={handleChange} />
                <input type="password" name="password" placeholder="Password" value={formData.password} onChange={handleChange} />
                <input type="text" name="path_to_crawler" placeholder="Path to Crawler" value={formData.path_to_crawler} onChange={handleChange} />
                <label>
                    <input type="checkbox" name="is_host" checked={formData.is_host} onChange={handleChange} />
                    Is Host
                </label>
                <input type="text" name="path_to_task_manager" placeholder="Path to Task Manager" value={formData.path_to_task_manager} onChange={handleChange} />
                <button onClick={handleSubmit}>Add Credentials</button>
            </div>

            <div style={{
                display: 'flex',
                flexDirection: 'row',
                
                justifyContent: 'center',
                alignItems: 'center',
                height: '30vh',
                width: '50vw',
                margin: 'auto',
                backgroundColor: 'lightgray',
                color: 'black',
            }}>
                {credentials.map((credential, index) => (
                    <div key={index} style={{
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'center',
                        alignItems: 'center',
                        height: '30vh',
                        width: '30%',
                        margin: 'auto',
                        backgroundColor: 'lightgray',
                        color: 'black',
                    }}>
                        <ul>
                            <li>IP: {credential.ip}</li>
                            <li>Username: {credential.username}</li>
                            <li>Path to Crawler: {credential.path_to_crawler}</li>
                            <li>Is Host: {credential.is_host.toString()}</li>
                            <li>Path to Task Manager: {credential.path_to_task_manager}</li>
                        </ul>
                        <button onClick={() => handleDelete(index)}>Delete</button>

                    </div>
                ))}
            </div>

            <div style={{
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: 'center',
                height: '30vh',
                width: '50vw',
                margin: 'auto',
                backgroundColor: 'lightgray',
                color: 'black',
            }}>
                <button onClick={() => 
                    runCrawlers(credentials)
                } style={{
                    width: '100%',
                    height: '100%',
                    margin: 'auto',
                    backgroundColor: 'lightblue',
                    color: 'black',
                
                }}>Run Crawler System</button>
            </div>
        </div>
    );
}

export default CrawlerRunner;
