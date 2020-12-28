import * as React from "react";
import { Admin, Resource} from 'react-admin';
import jsonServerProvider from 'ra-data-json-server';
import { ComputerList, ComputerEdit } from './models/computers'
import { IpxeaccountList, IpxeaccountCreate, IpxeaccountEdit } from './models/ipxeaccount'
import "./App.css"
import ComputerIcon from '@material-ui/icons/Computer';
import VpnKeyIcon from '@material-ui/icons/VpnKey';

const dataProvider = jsonServerProvider('/api/v1');
const App = () => (
    <Admin dataProvider={dataProvider}>
        <Resource name="computers" icon={ComputerIcon} list={ComputerList} edit={ComputerEdit}/>
        <Resource name="ipxeaccounts" options={{ label: 'iPXE accounts' }} icon={VpnKeyIcon} list={IpxeaccountList} edit={IpxeaccountEdit} create={IpxeaccountCreate}/>
    </Admin>
);

export default App;
