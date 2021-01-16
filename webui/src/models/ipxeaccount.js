import * as React from "react";
import {
    BooleanField,
    BooleanInput,
    Create,
    Datagrid,
    DateField,
    DateTimeInput,
    Edit,
    EditButton,
    Filter,
    List,
    Pagination,
    SimpleForm,
    TextField,
    TextInput,
    required,
} from 'react-admin';

const PostPagination = props => <Pagination rowsPerPageOptions={[15, 30, 50, 100, 200]} {...props} />;

const IpxeaccountFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Username" source="username" alwaysOn />
        <BooleanInput label="Is Admin" source="is_admin" alwaysOn />
    </Filter>
);

export const IpxeaccountList = props => (
    <List pagination={<PostPagination />} filters={<IpxeaccountFilter />} sort={{ field: 'username', order: 'ASC' }} {...props}>
        <Datagrid>
            <TextField source="username" />
            <BooleanField source="is_admin" />
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
            <TextInput source="password" type="password" validate={[required()]} />
            <TextInput source="password_confirmation" type="password" validate={[required()]} />
            <BooleanInput source="is_admin" helperText="Used if API auth is enabled" />
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
            <BooleanInput source="is_admin" helperText="Used if API auth is enabled" />
            <DateTimeInput source="last_login" disabled />
            <DateTimeInput source="created_at" disabled />
            <DateTimeInput source="updated_at" disabled />
        </SimpleForm>
    </Edit>
);