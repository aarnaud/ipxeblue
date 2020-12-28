import * as React from "react";
import {
    Create,
    Datagrid,
    DateField,
    DateTimeInput,
    Edit,
    EditButton,
    Filter,
    List,
    SimpleForm,
    TextField,
    TextInput,
    required,
} from 'react-admin';


const IpxeaccountFilter = (props) => (
    <Filter {...props}>
        <TextInput label="username" source="username" alwaysOn />
    </Filter>
);

export const IpxeaccountList = props => (
    <List filters={<IpxeaccountFilter />} {...props}>
        <Datagrid>
            <TextField source="username" />
            <DateField source="last_login" showTime={true} />
            <DateField source="created_at" showTime={true} />
            <DateField source="updated_at" showTime={true} />
            <EditButton />
        </Datagrid>
    </List>
);

export const IpxeaccountCreate = props => (
    // undoable={false} disable optimistic rendering
    <Create undoable={false} {...props}>
        <SimpleForm>
            <TextInput source="username" validate={[required()]} />
            <TextInput source="password" validate={[required()]} />
            <TextInput source="password_confirmation" validate={[required()]} />
        </SimpleForm>
    </Create>
);

export const IpxeaccountEdit = props => (
    // undoable={false} disable optimistic rendering
    <Edit undoable={false} {...props}>
        <SimpleForm>
            <TextInput source="username" disabled />
            <TextInput source="password" type="password" />
            <TextInput source="password_confirmation" type="password" />
            <DateTimeInput source="last_login" disabled />
            <DateTimeInput source="created_at" disabled />
            <DateTimeInput source="updated_at" disabled />
        </SimpleForm>
    </Edit>
);