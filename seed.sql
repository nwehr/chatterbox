-- drop role if exists chatterbox;
-- create role chatterbox with login password 'chatterbox';

drop table if exists message_recipients;
drop table if exists messages;
drop table if exists recipients;

create table recipients (
    "identity" varchar(128) not null
    , "password" varchar(64) not null
    , "public_keys" varchar(1024)[]
    , primary key("identity")
);

create table messages (
    "uuid" uuid not null
    , "parent_uuid" uuid default null
    , "conversation_id" varchar(64) not null
    , "recipients" varchar(128)[]
    , "from" varchar(128) not null
    , "encoding" varchar(128)
    , "length" bigint not null default 0
    , "data" text
    , "created_at" timestamp with time zone not null default now()
    , primary key("uuid")
    , constraint from_fkey foreign key("from") references recipients("identity")
    , constraint parent_uuid_fkey foreign key("parent_uuid") references messages("uuid")
);

create table message_recipients (
    "message_uuid" uuid
    , "conversation_id" varchar(64) not null
    , "recipient" varchar(128)
    , "read_at" timestamp with time zone default null
    , primary key("message_uuid", "recipient")
    , constraint message_uuid_fkey foreign key("message_uuid") references messages("uuid")
    , constraint recipient_fkey foreign key("recipient") references recipients("identity")
);

insert into recipients (
    "identity"
    , "password"
    , "public_keys"
) 
values 
	(
        '@nate.errorcode.io'
        , encode(sha256(('catch22' || 'jiemahGu0saoP3aiwieC8Eezeexeevai8aeSot9nah7xaV6vf')::bytea), 'hex')
        , '{age1rmpjlh40vsmry47pad0h4u0lavtrm0nlypaya4adf7xy9n0rd5zqzpgkua}'
    )
	, (
        '@kevpatt.errorcode.io'
        , encode(sha256(('catch22' || 'jiemahGu0saoP3aiwieC8Eezeexeevai8aeSot9nah7xaV6vf')::bytea), 'hex')
        , '{age1qz33qeyel6uzyvl9kfscmx4le7tkylqecn2lhc820r7gkj5dz9zqdycvkw}'
    )
;

-- with new_messages as (
--     insert into messages (
--         "uuid"
--         , "conversation_id"
--         , "recipients"
--         , "from"
--         , "encoding"
--         , "data"
--         , "created_at"
--     ) 
--     values (
--         gen_random_uuid()
--         , encode(sha256('@kevpatt.errorcode.io;@nate.errorcode.io'), 'hex')
--         , '{@kevpatt.errorcode.io, @nate.errorcode.io}'
--         , '@kevpatt.errorcode.io'
--         , 'text/plain'
--         , 'Hey! You there?'
--         , now()
--     )
--     , (
--         gen_random_uuid()
--         , encode(sha256('@kevpatt.errorcode.io;@nate.errorcode.io'), 'hex')
--         , '{@kevpatt.errorcode.io, @nate.errorcode.io}'
--         , '@nate.errorcode.io'
--         , 'text/plain'
--         , 'Sure am! Whats up?'
--         , now()
--     )
--     , (
--         gen_random_uuid()
--         , encode(sha256('@kevpatt.errorcode.io;@nate.errorcode.io'), 'hex')
--         , '{@kevpatt.errorcode.io, @nate.errorcode.io}'
--         , '@kevpatt.errorcode.io'
--         , 'text/plain'
--         , 'Have you seen the new Macbook Pros?'
--         , now()
--     )
--     , (
--         gen_random_uuid()
--         , encode(sha256('@nate.errorcode.io'), 'hex')
--         , '{@nate.errorcode.io}'
--         , '@nate.errorcode.io'
--         , 'text/plain'
--         , 'Just talking to myself'
--         , now()
--     ) 
--     returning "uuid", "recipients", "conversation_id"
-- )
-- insert into message_recipients ("message_uuid", "recipient", "conversation_id") select "uuid", unnest("recipients"), "conversation_id" from new_messages;

grant all privileges on all tables in schema public to chatterbox;