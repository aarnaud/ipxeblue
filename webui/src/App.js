import * as React from "react";
import { Admin, Resource} from 'react-admin';
import jsonServerProvider from 'ra-data-json-server';
import { ComputerList, ComputerEdit } from './models/computers'
import { IpxeaccountList, IpxeaccountCreate, IpxeaccountEdit } from './models/ipxeaccount'
import { BootentryList, BootentryCreate, BootentryEdit } from './models/bootentry'
import "./App.css"
import ComputerIcon from '@material-ui/icons/Computer';
import VpnKeyIcon from '@material-ui/icons/VpnKey';
import AssignmentIcon from '@material-ui/icons/Assignment';

import { createMuiTheme } from '@material-ui/core/styles';
import blueGrey from '@material-ui/core/colors/blueGrey';
import lightBlue from '@material-ui/core/colors/lightBlue';
const theme = createMuiTheme({
    typography: {
        fontSize: 12,
    },
    palette: {
        primary: lightBlue,
        secondary: blueGrey,
        type: 'dark', // Switching the dark mode on is a single property value change.
    },
});


const apiUrl = '/api/v1';
const dataProvider = jsonServerProvider(apiUrl);
const myDataProvider = {
    ...dataProvider,
    update: (resource, params) => {
        if (resource !== 'bootentries' || !params.data.files) {
            // fallback to the default implementation
            return dataProvider.update(resource, params);
        }

        // remove file if fileobject is null
        params.data.files = params.data.files.filter(file => {
            if (file === undefined){
                return false
            }
            return file.file !== null
        })

        // set name from fileobject title
        params.data.files.map(file => {
            if (file.file) {
                file.name = file.file.title
            }
            return file
        })

        // if rawfile exist it's new file, need to be upload
        const newFiles = params.data.files.filter(
            file => {
                return file.file.rawFile instanceof File
            }
        );
        return Promise.all(newFiles.map(file => {
            return fileUpload(file.file, `${apiUrl}/${resource}/${params.id}/files/${file.name}`)
        })).then(files =>
            dataProvider.update(resource, params)
        )
    },
};

const fileUpload = (file, url) => {
    const formData = new FormData()
    formData.append('file', file.rawFile)
    return  fetch(url, {
        method: 'POST',
        body: formData
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
        })
        .catch(error => {
            console.error(error)
        })
}



const App = () => (
    <Admin dataProvider={myDataProvider} theme={theme}>
        <Resource name="computers" icon={ComputerIcon} list={ComputerList} edit={ComputerEdit}/>
        <Resource name="ipxeaccounts" options={{ label: 'iPXE accounts' }} icon={VpnKeyIcon} list={IpxeaccountList} edit={IpxeaccountEdit} create={IpxeaccountCreate}/>
        <Resource name="bootentries" options={{ label: 'Boot entries' }} icon={AssignmentIcon} list={BootentryList} edit={BootentryEdit} create={BootentryCreate}/>
    </Admin>
);

export default App;
