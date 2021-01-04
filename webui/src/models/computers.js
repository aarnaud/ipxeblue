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
    SelectInput,
    SimpleForm,
    SimpleFormIterator,
    ReferenceInput,
    ReferenceField,
    TextField,
    TextInput,
    required,
} from 'react-admin';


const ComputerFilter = (props) => (
    <Filter {...props}>
        <TextInput label="name" source="name" alwaysOn />
        <TextInput label="MAC" source="mac" alwaysOn />
        <TextInput label="IP" source="ip" alwaysOn />
        <TextInput label="Manufacturer" source="manufacturer" alwaysOn />
        <TextInput label="Product" source="product" alwaysOn />
        <TextInput label="Serial" source="serial" alwaysOn />
        <TextInput label="Build Arch" source="build_arch" allowEmpty />
        <TextInput label="Platform" source="platform" allowEmpty />
        <ReferenceInput label="Bootentry" source="bootentry_uuid" alwaysOn allowEmpty={true} reference="bootentries">
            <SelectInput optionText="name" />
        </ReferenceInput>
    </Filter>
);

export const ComputerList = props => (
    <List filters={<ComputerFilter />} {...props}>
        <Datagrid>
            <TextField source="name" />
            <TextField source="id" />
            <TextField source="mac" />
            <TextField source="ip" />
            <TextField source="hostname" />
            <DateField source="last_seen" showTime={true} />
            <TextField source="platform" />
            <TextField source="build_arch" />
            <TextField source="manufacturer" />
            <TextField source="product" />
            <TextField source="serial" />
            <TextField source="version" />
            <ReferenceField label="Bootentry" source="bootentry_uuid" reference="bootentries">
                <TextField source="name" />
            </ReferenceField>
            <EditButton />
        </Datagrid>
    </List>
);

export const ComputerEdit = props => (
    // undoable={false} disable optimistic rendering
    <Edit undoable={false} {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <ReferenceInput label="Bootentry" source="bootentry_uuid" allowEmpty={true} reference="bootentries">
                <SelectInput optionText="description" />
            </ReferenceInput>
            <TextInput source="id" disabled />
            <TextInput source="mac" disabled />
            <TextInput source="ip" disabled />
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