import * as React from "react";
import {
    ArrayInput,
    Datagrid,
    DateField,
    DateTimeInput,
    Edit,
    EditButton,
    Filter,
    List,
    SimpleForm,
    SimpleFormIterator,
    TextField,
    TextInput,
    required,
} from 'react-admin';


const ComputerFilter = (props) => (
    <Filter {...props}>
        <TextInput label="name" source="name" alwaysOn />
        <TextInput label="MAC" source="mac" allowEmpty />
        <TextInput label="Build Arch" source="build_arch" allowEmpty />
        <TextInput label="Platform" source="platform" allowEmpty />
    </Filter>
);

export const ComputerList = props => (
    <List filters={<ComputerFilter />} {...props}>
        <Datagrid>
            <TextField source="name" />
            <TextField source="id" />
            <TextField source="mac" />
            <TextField source="hostname" />
            <DateField source="last_seen" showTime={true} />
            <TextField source="platform" />
            <TextField source="build_arch" />
            <TextField source="manufacturer" />
            <TextField source="product" />
            <TextField source="serial" />
            <TextField source="asset" />
            <TextField source="version" />
            <EditButton />
        </Datagrid>
    </List>
);

export const ComputerEdit = props => (
    // undoable={false} disable optimistic rendering
    <Edit undoable={false} {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <TextInput source="id" disabled />
            <TextInput source="mac" disabled />
            <TextInput source="hostname" disabled />
            <DateTimeInput source="last_seen" disabled />
            <TextInput source="platform" disabled />
            <TextInput source="build_arch" disabled />
            <TextInput source="manufacturer" disabled />
            <TextInput source="product" disabled />
            <TextInput source="serial" disabled />
            <TextInput source="asset" disabled />
            <TextInput source="version" disabled />
            <DateTimeInput source="last_login" disabled />
            <DateTimeInput source="created_at" disabled />
            <DateTimeInput source="updated_at" disabled />
            <ArrayInput source="tags">
                <SimpleFormIterator>
                    <TextInput source="key" label="key" validate={[required()]} />
                    <TextInput source="value" label="value" />
                </SimpleFormIterator>
            </ArrayInput>
        </SimpleForm>
    </Edit>
);