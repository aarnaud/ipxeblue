import * as React from "react";
import { Admin, Resource} from 'react-admin';
import jsonServerProvider from 'ra-data-json-server';
import { ComputerList, ComputerEdit } from './models/computers'
import "./App.css"
import ComputerIcon from '@material-ui/icons/Computer';

const dataProvider = jsonServerProvider('/api/v1');
const App = () => (
    <Admin dataProvider={dataProvider}>
        <Resource name="computers" icon={ComputerIcon} list={ComputerList} edit={ComputerEdit}/>
    </Admin>
);

export default App;
