import * as React from "react";
import {
    ArrayField,
    ArrayInput,
    ChipField,
    Datagrid,
    DateField,
    DateTimeInput,
    DeleteButton,
    Edit,
    EditButton,
    Filter,
    List,
    Pagination,
    ReferenceInput,
    SelectInput,
    SimpleForm,
    SimpleFormIterator,
    SingleFieldList,
    TextField,
    TextInput,
    required,
} from 'react-admin';

const PostPagination = props => <Pagination rowsPerPageOptions={[15, 30, 50, 100, 200]} {...props} />;
const EditTitle = ({ record }) => {
    return <span>{record ? `${record.name}` : ''}</span>;
};

const ComputerFilter = (props) => (
    <Filter {...props}>
        <TextInput label="name" source="name" alwaysOn />
        <TextInput label="MAC" source="mac" alwaysOn />
        <TextInput label="IP" source="ip" alwaysOn />
        <TextInput label="Serial" source="serial" alwaysOn />
        <ReferenceInput label="Bootentry" source="bootentry_uuid" alwaysOn allowEmpty={true} reference="bootentries">
            <SelectInput optionText="name" />
        </ReferenceInput>
        <TextInput label="Manufacturer" source="manufacturer" allowEmpty />
        <TextInput label="Product" source="product" allowEmpty />
        <TextInput label="Build Arch" source="build_arch" allowEmpty />
        <TextInput label="Platform" source="platform" allowEmpty />
        <TextInput label="Tags" source="value" allowEmpty />
    </Filter>
);

export const ComputerList = props => (
    <List pagination={<PostPagination />} perPage={15} filters={<ComputerFilter />} sort={{ field: 'name', order: 'ASC' }} {...props}>
        <Datagrid>
            <TextField source="name" />
            <TextField source="mac" />
            <TextField source="ip" />
            <DateField source="last_seen" showTime={true} />
            <TextField source="manufacturer" />
            <TextField source="product" />
            <TextField source="serial" />
            <TextField source="version" />
            <ArrayField source="bootorder">
                <SingleFieldList linkType={false}>
                    <ChipField source="name" label="Name"  />
                </SingleFieldList>
            </ArrayField>
            <ArrayField source="tags">
                <SingleFieldList linkType={false}>
                    <ChipField source="value" label="value"  />
                </SingleFieldList>
            </ArrayField>
            <EditButton />
            <DeleteButton />
        </Datagrid>
    </List>
);

export const ComputerEdit = props => (
    <Edit mutationMode="pessimistic" title={<EditTitle />} {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <ArrayInput source="bootorder">
                <SimpleFormIterator>
                    <ReferenceInput label="Bootentry" source="id" validate={[required()]} reference="bootentries" >
                        <SelectInput optionText="description" />
                    </ReferenceInput>
                </SimpleFormIterator>
            </ArrayInput>
            <ArrayInput source="tags">
                <SimpleFormIterator>
                    <TextInput source="key" label="key" validate={[required()]} />
                    <TextInput source="value" label="value" />
                </SimpleFormIterator>
            </ArrayInput>
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
            <ReferenceInput label="Last iPXE account logged" source="last_ipxeaccount" allowEmpty={true} reference="ipxeaccounts" disabled>
                <SelectInput optionText="username" />
            </ReferenceInput>
            <DateTimeInput source="created_at" disabled />
            <DateTimeInput source="updated_at" disabled />
        </SimpleForm>
    </Edit>
);