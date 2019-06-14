create or replace function password_validation( password text ) returns text
as $$
  begin

    if length( password ) < 6 then
      return 'Password must be at least 6 characters long' ;
    end if;

    return null;

  end;
$$ language 'plpgsql';


create or replace function login_validation( login text ) returns text
as $$
  declare
    i integer;
    array_elem integer;
    login_array char[] := string_to_array( login , null ) ;
  begin

    if length( login ) < 3 then
      return 'Login must be at least 3 characters long' ;
    end if;

    for i in 1..array_length(login_array, 1)
      loop
        array_elem = ascii(login_array[i]) ;
        if array_elem != 95 and not ( array_elem >= 48 and array_elem <= 57 ) and
           not ( array_elem >= 65 and array_elem <= 90 ) and
           not ( array_elem >= 97 and array_elem <= 122 )
          then
          return 'Login can contain only Latin characters, numbers and underscore' ;
        end if;
      end loop;

    return null;

  end;
$$ language 'plpgsql';
