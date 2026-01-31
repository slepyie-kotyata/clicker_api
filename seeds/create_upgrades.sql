BEGIN;

--0 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (1, 'Сэндвич', 'sandwich', 'dish', 1, 0, 0);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (1, 'mPc', 1, 1);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (2, 'Газировка', 'cola', 'dish', 1, 25, 0);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (2, 'mPc', 1, 2);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (3, 'Картошка фри', 'french_fries', 'dish', 1, 40, 0);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (3, 'mPc', 1, 3);

--1 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (4, 'Новая плита', 'new_stove', 'equipment', 3.2, 250, 1);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (4, 'dPc', 1, 4);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (5, 'Наггетсы', 'nuggets', 'dish', 1, 60, 1);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (5, 'mPc', 1, 5);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (6, 'Гамбургер', 'hamburger', 'dish', 1, 90, 1);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (6, 'mPc', 1, 6);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (7, 'Картоф. дольки', 'potato_wedges', 'dish', 1, 130, 1);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (7, 'mPc', 1, 7);

--2 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (8, 'Хот-дог', 'hot_dog', 'dish', 1, 200, 2);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (8, 'mPc', 1, 8);

--3 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (9, 'Чизбургер', 'cheeseburger', 'dish', 1, 320, 3);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (9, 'mPc', 1, 9);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (10, 'Газовая горелка', 'gas_burner', 'equipment', 3.5, 800, 3);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (10, 'dPc', 2, 10);

--4 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (11, 'Милкшейк клубника', 'strawberry_milkshake', 'dish', 1, 500, 4);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (11, 'mPc', 1, 11);

--5 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (12, 'Тако', 'taco', 'dish', 1, 800, 5);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (12, 'mPc', 1, 12);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (13, 'Панини', 'panini', 'dish', 1, 1200, 5);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (13, 'mPc', 1, 13);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (14, 'Вафля', 'waffle', 'dish', 1, 1800, 5);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (14, 'mPc', 2, 14);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (15, 'Доп. касса', 'new_cashier', 'equipment', 3, 1500, 5);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (15, 'sPs', 2, 15);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (16, 'Проф. духовка', 'pro_oven', 'equipment', 5, 4000, 5);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (16, 'dPs', 1, 16);

--6 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (17, 'Буррито', 'burrito', 'dish', 1, 2600, 6);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (17, 'mPc', 2, 17);

--7 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (18, 'Багет', 'ham_baguette', 'dish', 1, 3800, 7);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (18, 'mPc', 2, 18);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (19, 'Умная касса', 'auto_cashier', 'equipment', 2.3, 15000, 7);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (19, 'mPs', 1, 19);

--8 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (20, 'Пончик', 'donate', 'dish', 1, 5500, 8);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (20, 'mPc', 3, 20);

--9 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (21, 'Мороженое', 'ice_cream', 'dish', 1, 8000, 9);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (21, 'mPc', 3, 21);

--10 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (22, 'Тирамису', 'tiramisu', 'dish', 1, 12000, 10);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (22, 'mPc', 4, 22);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (23, 'Милкшейк ваниль', 'vanilla_milkshake', 'dish', 1, 18000, 10);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (23, 'mPc', 4, 23);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (24, 'Круассаны', 'croissants', 'dish', 1, 26000, 10);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (24, 'mPc', 5, 24);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (25, 'Реклама блюд', 'promo_dishes', 'global', 3.5, 15000, 10);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (25, 'dM', 2, 25);

--11 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (26, 'Чизкейк', 'cheesecake', 'dish', 1, 26000, 11);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (26, 'mPc', 6, 26);

--12 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (27, 'Сок', 'juice', 'dish', 1, 55000, 12);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (27, 'mPc', 7, 27);

--13 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (28, 'Кексы', 'cupcakes', 'dish', 1, 55000, 13);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (28, 'mPc', 8, 28);

--14 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (29, 'Панкейки', 'pancakes', 'dish', 1, 120000, 14);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (29, 'mPc', 9, 29);

--15 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (30, 'Американо', 'coffee', 'dish', 1, 180000, 15);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (30, 'mPc', 10, 30);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (31, 'Макаруны', 'macaroni', 'dish', 1, 180000, 15);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (31, 'mPc', 11, 31);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (32, 'Капучино', 'cappuccino', 'dish', 1, 380000, 15);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (32, 'mPc', 12, 32);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (33, 'Повышение цен', 'price_increase', 'global', 3, 50000, 15);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (33, 'mM', 1.5, 33);

--16 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (34, 'Морской салат', 'sea_salad', 'dish', 1, 550000, 15);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (34, 'mPc', 13, 34);

--17 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (35, 'Чай', 'tea', 'dish', 1, 800000, 17);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (35, 'mPc', 15, 35);

--18 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (36, 'Омлет', 'omelette', 'dish', 1, 1200000, 18);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (36, 'mPc', 16, 36);

--19 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (37, 'Клубничный смузи', 'strawberry_smoothie', 'dish', 1, 1800000, 19);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (37, 'mPc', 19, 37);

--20 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (38, 'Пицца', 'pizza', 'dish', 1, 3000000, 20);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (38, 'mPc', 21, 38);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (39, 'Ролл с креветкой', 'shrimp_roll', 'dish', 1, 4500000, 20);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (39, 'mPc', 25, 39);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (40, 'Палтус', 'grilled_halibut', 'dish', 1, 6500000, 20);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (40, 'mPc', 28, 40);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (41, 'Современный гриль', 'next_gen_grill', 'equipment', 2.8, 200000, 20);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (41, 'dM', 2, 41);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (42, 'Качество персонала', 'staff_quality', 'global', 3, 95000, 20);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (42, 'sPs', 10, 42);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (43, 'Официант', 'waiter', 'staff', 2.3, 100000, 20);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (43, 'mPs', 10, 43);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (44, 'Шеф-повар', 'chef', 'staff', 2.3, 80000, 20);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (44, 'dPs', 10, 44);

-- 21 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (45, 'Мохито', 'mojito', 'dish', 1, 9500000, 21);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (45, 'mPc', 32, 45);

-- 22 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (46, 'Спагетти', 'spaghetti', 'dish', 1, 14000000, 22);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (46, 'mPc', 36, 46);

-- 23 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (47, 'Суфле', 'caramel_souffle', 'dish', 1, 20000000, 23);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (47, 'mPc', 17, 47);

-- 24 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (48, 'Вино', 'wine', 'dish', 1, 30000000, 24);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (48, 'mPc', 20, 48);

-- 25 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (49, 'Рамен', 'ramen', 'dish', 1, 45000000, 25);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (49, 'mPc', 23, 49);

-- 26 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (50, 'Стейк', 'steak', 'dish', 1, 70000000, 26);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (50, 'mPc', 26, 50);

-- 27 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (51, 'Дракон-ролл', 'dragon_roll', 'dish', 1, 110000000, 27);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (51, 'mPc', 32, 51);

-- 28 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (52, 'Том-Ям', 'tom_yam', 'dish', 1, 170000000, 28);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (52, 'mPc', 38, 52);

-- 29 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (53, 'Суши-сет', 'sushi_set', 'dish', 1, 260000000, 29);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (53, 'mPc', 46, 53);

-- 30 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (54, 'Авторское', 'signature_dish', 'dish', 1, 400000000, 30);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (54, 'mPc', 56, 54);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (55, 'Апероль Спритц', 'aperol_spritz', 'dish', 1, 600000000, 30);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (55, 'mPc', 67, 55);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (56, 'Мясная нарезка', 'beef_carpaccio', 'dish', 1, 900000000, 30);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (56, 'mPc', 79, 56);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (57, 'Авто-продажи', 'auto_sales', 'global', 2.4, 500000, 30);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (57, 'mpM', 2, 57);

-- 35 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (58, 'Френч-сет', 'french_set', 'dish', 1, 1400000000, 35);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (58, 'mPc', 95, 58);

-- 40 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (59, 'Рыбный микс', 'fish_mix', 'dish', 1, 2200000000, 40);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (59, 'mPc', 115, 59);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (60, 'Блюдо от шефа', 'chef_special', 'dish', 1, 3500000000, 40);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (60, 'mPc', 139, 60);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (61, 'Бабл-ти матча', 'matcha_boba_ice', 'dish', 1, 5500000000, 40);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (61, 'mPc', 170, 61);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (62, 'Экспресс-доставка', 'express_delivery', 'global', 2.8, 5000000, 40);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (62, 'mpM', 5, 62);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (63, 'Быстрый расчет', 'auto_cash', 'global', 2, 10000000, 40);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (63, 'sPs', 100, 63);

-- 45 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (64, 'Тропический сет', 'tropical_surf_set', 'dish', 1, 8500000000, 45);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (64, 'mPc', 206, 64);

-- 50 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (65, 'Конфи из утки', 'duck_confit', 'dish', 1, 13000000000, 50);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (65, 'mPc', 251, 65);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (66, 'Бургер с трюфелем', 'truffle_burger', 'dish', 1, 20000000000, 50);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (66, 'mPc', 305, 66);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (67, 'Бабл-ти анчан', 'blue_matcha_boba_ice', 'dish', 1, 32000000000, 50);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (67, 'mPc', 371, 67);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (68, 'Фуа-гра сет', 'imperial_foie_gras', 'dish', 1, 50000000000, 50);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (68, 'mPc', 450, 68);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (69, 'Мишлен-фантазия', 'michelin_set', 'dish', 1, 80000000000, 50);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (69, 'mPc', 547, 69);

-- 55 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (70, 'Морская сет', 'sea_plate', 'dish', 1, 125000000000, 55);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (70, 'mPc', 666, 70);

-- 60 уровень
INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (71, 'Матча латте', 'matcha_latte', 'dish', 1, 200000000000, 60);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (71, 'mPc', 812, 71);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (72, 'Супрем-блюдо', 'supreme_maestro_dish', 'dish', 1, 320000000000, 60);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (72, 'mPc', 984, 72);

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (73, 'Десертный сет', 'dessert_set', 'dish', 1, 520000000000, 60);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (73, 'mPc', 1192, 73);

-- 65 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (74, 'Панна-котта', 'jiggly_cat', 'dish', 1, 900000000000, 65);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (74, 'mPc', 1455, 74);

-- 70 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (75, 'Фудтрак', 'food_truck', 'point', 2.2, 100000000000, 70);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (75, 'mPs', 500, 75);

-- 75 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (76, 'Маленькое кафе', 'small_cafe', 'point', 2.2, 150000000000, 75);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (76, 'mPs', 1500, 76);

-- 80 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (77, 'Ресторан', 'family_restaurant', 'point', 2.3, 500000000000, 80);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (77, 'mPs', 5000, 77);

-- 90 уровень

INSERT INTO public.upgrades (id, name, icon_name, upgrade_type, price_factor, price, access_level) VALUES (78, 'Гастро-ресторан', 'gastro_restaurant', 'point', 2.4, 1500000000000, 90);
INSERT INTO public.boosts (id, boost_type, value, upgrade_id) VALUES (78, 'mPs', 15000, 78);

COMMIT;
