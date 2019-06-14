-- курсы
insert into courses values (default,'Python');
-- раздел
insert into sections values (
                                default,
                                (select c.id from courses c where c.name='Python') ,
                                'Основи програмування на Python'
                            );
-- урок
insert into lessons values (
                               default,
                               (select s.id from sections s where s.name='Основи програмування на Python') ,
                               'Історія мови програмування'
                           );
-- параграф
insert into paragraphs_or_tests values (
                               default,
                               (select l.id from lessons l where l.name='Історія мови програмування') ,
                               'Створення Python',
                               'Розробка мови Python була розпочата в кінці 1980-х років голландським програмістом Гвідо ван Россумом',
                               0
                           );
-- тест
insert into paragraphs_or_tests values (
                                           default,
                                           (select l.id from lessons l where l.name='Історія мови програмування') ,
                                           'Створення Python',
                                           'Коли було розпочато розробку мови програмування Python?',
                                           10
                                       );
-- варианты ответов теста
insert into tests_answers values (
                                     default,
                                     (select t.id from paragraphs_or_tests t where t.text='Коли було розпочато розробку мови програмування Python?') ,
                                     'кінець 1980-х',
                                     true
                                 );
-- варианты ответов теста
insert into tests_answers values (
                                     default,
                                     (select t.id from paragraphs_or_tests t where t.text='Коли було розпочато розробку мови програмування Python?') ,
                                     'кінець 1990-х',
                                     false
                                 );
-- варианты ответов теста
insert into tests_answers values (
                                     default,
                                     (select t.id from paragraphs_or_tests t where t.text='Коли було розпочато розробку мови програмування Python?') ,
                                     'кінець 1970-х',
                                     false
                                 );
-- варианты ответов теста
insert into tests_answers values (
                                     default,
                                     (select t.id from paragraphs_or_tests t where t.text='Коли було розпочато розробку мови програмування Python?') ,
                                     'початок 2000-х',
                                     false
                                 );


