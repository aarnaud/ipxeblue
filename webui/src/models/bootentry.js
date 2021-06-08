import * as React from "react";
import {
    ArrayInput,
    BooleanField,
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
    Pagination,
    SimpleForm,
    SimpleFormIterator,
    TextField,
    TextInput,
    required,
} from 'react-admin';

const PostPagination = props => <Pagination rowsPerPageOptions={[15, 30, 50, 100, 200]} {...props} />;
const EditTitle = ({ record }) => {
    return <span>{record ? `${record.name}` : ''}</span>;
};

const BootentryFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Name" source="name" alwaysOn />
        <TextInput label="Description" source="description" alwaysOn />
        <BooleanInput label="Persistent" source="persistent" alwaysOn />
    </Filter>
);

export const BootentryList = props => (
    <List pagination={<PostPagination />} perPage={15} filters={<BootentryFilter />} sort={{ field: 'name', order: 'ASC' }} {...props}>
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
    <Edit mutationMode="pessimistic" title={<EditTitle />} {...props}>
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
                    <TextInput source="subpath" label="Sub Path" />
                    <FileInput source="file" label="File" placeholder="click to select file to upload" options={{noDrag: true, maxFiles: 1, multiple: false}}>
                        <FileField source="src" title="title" />
                    </FileInput>
                </SimpleFormIterator>
            </ArrayInput>
        </SimpleForm>
    </Edit>
);