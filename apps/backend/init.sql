-- Create profiles table
create table public.profiles (
  user_id text primary key,
  data_blob text not null,
  updated_at timestamp with time zone default timezone('utc'::text, now()) not null
);

-- Turn on Row Level Security
alter table public.profiles enable row level security;

-- For MVP: Allow public access (or create specific policies)
-- Since we are using a "demo-user-001" and just an API Key, we permit interaction.
create policy "Public profiles are viewable by everyone."
  on profiles for select
  using ( true );

create policy "Users can insert their own profile."
  on profiles for insert
  with check ( true );

create policy "Users can update their own profile."
  on profiles for update
  using ( true );
