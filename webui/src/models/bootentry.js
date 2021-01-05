import * as React from "react";
import {
    BooleanInput,
    Create,
    Datagrid,
    DateField,
    DateTimeInput,
    Edit,
    EditButton,
    FileField,
    FileInput,
    Filter,
    List,
    SimpleForm,
    TextField,
    TextInput,
    required, SimpleFormIterator, ArrayInput, BooleanField,
} from 'react-admin';

const BootentryFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Name" source="name" alwaysOn />
        <TextInput label="Description" source="description" alwaysOn />
        <BooleanInput label="Description" source="description" alwaysOn />
    </Filter>
);

export const BootentryList = props => (
    <List filters={<BootentryFilter />} {...props}>
        <Datagrid>
            <TextField source="id" disabled />
            <TextField source="name" />
            <TextField source="description" />
            <BooleanField source="persistent" />
            <DateField source="created_at" showTime={true} />
            <DateField source="updated_at" showTime={true} />
            <EditButton />
        </Datagrid>
    </List>
);

export const BootentryCreate = props => (
    // undoable={false} disable optimistic rendering
    <Create undoable={false} {...props}>
        <SimpleForm>
            <TextInput source="name" validate={[required()]} />
            <TextInput source="description" validate={[required()]} />
            <BooleanInput source="persistent" />
            <TextInput component="pre" fullWidth={true} rows={10} options={{ multiline: true }} source="ipxe_script" />
        </SimpleForm>
    </Create>
);

export const BootentryEdit = props => (
    // undoable={false} disable optimistic rendering
    <Edit undoable={false} {...props}>
        <SimpleForm>
            <TextInput source="id" disabled />
            <DateTimeInput source="created_at" disabled />
            <DateTimeInput source="updated_at" disabled />
            <TextInput source="name" validate={[required()]} />
            <TextInput source="description" validate={[required()]} />
            <BooleanInput source="persistent" />
            <TextInput component="pre" fullWidth={true} rows={10} options={{ multiline: true }} source="ipxe_script" />
            <ArrayInput source="files">
                <SimpleFormIterator>
                    <BooleanInput source="protected" label="Protected" />
                    <BooleanInput source="templatized" label="Templatized" />
                    <FileInput source="file" label="File" placeholder="click to select file to upload" options={{noDrag: true, maxFiles: 1, multiple: false}}>
                        <FileField source="src" title="title" />
                    </FileInput>
                </SimpleFormIterator>
            </ArrayInput>
        </SimpleForm>
    </Edit>
);