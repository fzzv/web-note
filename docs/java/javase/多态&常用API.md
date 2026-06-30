# 多态、内部类、常用API

## 1 多态

### 1.1 面向对象三大特征 ?

- 封装 , 继承 , 多态

### 1.2 什么是多态 ?

- 一个对象在不同时刻体现出来的不同形态

- 举例 : 一只猫对象
  - 我们可以说猫就是猫 : Cat cat = new Cat();
  - 我们也可以说猫是动物 : Animal cat = new Cat();
  - 这里对象在不同时刻,体现出来的不同形态 , 我们就可以理解为多态

### 1.3 多态的前提

- 有继承/实现关系
- 有方法重写 
- 父类的引用指向子类的对象

```java
package com.fan.polymorphism_demo;
/*
    多态的三个前提条件
        1 需要有继承/实现关系
        2 需要有方法重写
        3 父类的引用指向子类的对象
 */
public class AnimalTest {
    public static void main(String[] args) {
        // 3 父类的引用指向子类的对象
        // 多态形式对象
        Animal a = new Cat();
    }
}

class Animal{
    public void eat(){
        System.out.println("吃东西");
    }
}

class Cat extends Animal{
    @Override
    public void eat() {
        System.out.println("猫吃鱼....");
    }
}

```

### 1.4 多态的成员访问特点

- 构造方法 : 和继承一样 , 子类通过super()访问父类的构造方法
- 成员变量 : 编译看左边(父类) , 执行看左边(父类)
- 成员方法 : 编译看左边(父类) , 执行看右边(子类)

```java
package com.fan.polymorphism_demo;

/*
    多态的成员访问特点 :
        1 构造方法 : 和继承一样 , 都是通过super()访问父类的构造方法
        2 成员变量 : 编译看左边(父类) , 执行看左边(父类)
        3 成员方法 : 编译看左边(父类) , 执行看右边(子类) , 注意 , 如果执行时
            1) 子类没有回动态去找父类中的方法
            2) 子类的特有方法无法进行调用(多态的缺点)
 */
public class MemberTest {
    public static void main(String[] args) {
        // 父类的引用指向子类的对象
        Fu f = new Zi();

        // 多态对象调用成员变量
        System.out.println(f.num);

        // 多态对新乡调用调用成员方法
        f.show();

        // 多态对象不能调用子类特有的方法
        // f.show2();
    }
}

class Fu {
    int num = 100;

    public void show() {
        System.out.println("父类的show方法");
    }
}

class Zi extends Fu {
    int num = 10;

    public void show() {
        System.out.println("子类的show方法");
    }

    public void show2(){
        System.out.println("子类特有的方法");
    }
}

```

### 1.5 多态的优缺点

- 优点 : 提高代码的扩展性
- 缺点 : 不能调用子类特有的功能

```java
package com.fan.polymorphism_test;

public abstract class Animal {
    private String breed;
    private String color;

    public Animal() {
    }

    public Animal(String breed, String color) {
        this.breed = breed;
        this.color = color;
    }

    public String getBreed() {
        return breed;
    }

    public void setBreed(String breed) {
        this.breed = breed;
    }

    public String getColor() {
        return color;
    }

    public void setColor(String color) {
        this.color = color;
    }

    public abstract void eat();
}

```

```java
package com.fan.polymorphism_test;

public class Cat extends Animal {
    public Cat() {
    }

    public Cat(String breed, String color) {
        super(breed, color);
    }

    @Override
    public void eat() {
        System.out.println("猫吃鱼...");
    }

    public void catchMouse() {
        System.out.println("抓老鼠...");
    }
}

```

```java
package com.fan.polymorphism_test;

public class Dog extends Animal {
    public Dog() {
    }

    public Dog(String breed, String color) {
        super(breed, color);
    }

    @Override
    public void eat() {
        System.out.println("狗吃骨头!");
    }

    public void lookDoor(){
        System.out.println("狗看门...");
    }
}

```

```java
package com.fan.polymorphism_test;

public class Pig extends Animal {
    public Pig() {
    }

    public Pig(String breed, String color) {
        super(breed, color);
    }

    @Override
    public void eat() {
        System.out.println("猪拱白菜...");
    }

    public void sleep() {
        System.out.println("一直再睡...");
    }
}

```

```java
package com.fan.polymorphism_test;

/*
    如果方法的参数是一个类的话 , 那么调用此方法需要传入此类的对象 , 或者子类对象

    多态的好处 :
        提高代码的扩展性 , 灵活性
    多态的缺点:
        不能调用子类的特有功能
 */
public class AnimalTest {
    public static void main(String[] args) {
        useAnimal(new Cat());

        System.out.println("---------");

        useAnimal(new Dog());

        System.out.println("---------");

        useAnimal(new Pig());
    }

    public static void useAnimal(Animal a){// Animal a = new Dog()
        a.eat();
        // 多态不能访问子类特有的功能
        // 如果解决 ?

        // 向下转型
        if(a instanceof Cat) {
            Cat cat = (Cat) a;
            cat.catchMouse();
        }
        if(a instanceof Dog) {
            Dog dog = (Dog) a;
            dog.lookDoor();
        }

        if(a instanceof Pig) {
            ((Pig) a).sleep();
        }
    }

//    // 定义一个使用猫类的方法
//    public static void useAnimal(Cat c) {// Cat c = new Cat();
//        c.eat();
//        c.catchMouse();
//    }
//
//    // 定义一个使用狗类的方法
//    public static void useAnimal(Dog d) {// Dog d = new Dog();
//        d.eat();
//        d.lookDoor();
//    }
//
//    // 定义一个使用猪类的方法
//    public static void useAnimal(Pig pig) {
//        pig.eat();
//        pig.sleep();
//    }
}

```

### 1.6 多态的转型

- 向上转型 :  把子类类型数据转成父类类型数据  Animal a = new Cat();
- 向下转型 :  把父类类型数据转成子类类型数据  Cat cat = (Cat)a;

### 1.7 多态的转型注意

- 如果被转的对象 , 对应的实际类型和目标类型不是同一种数据类型 , 那么转换时会出现ClassCastException异常

- ```java
  异常代码如下
  public static void main(String[] args) {
      Animal a = new Cat();
      useAnimal(a);
  }
  public static void useAnimal(Animal a) {
      Dog d = (Dog) a;
      d.eat();
  }
  ```

### 1.8 解决转型安全隐患

- 使用关键字   instanceof
- 作用 : 判断一个对象是否属于一种引用数据类型
- 格式 : 对象名 instanceof 引用数据类型
  - 通俗的理解：判断关键字左边的变量，是否是右边的类型，返回boolean类型结果

## 2 内部类

### 2.1 内部类的分类

- ##### 什么是内部类 ?

  - 一个A类 中 定义一个B类 , 那么B类就属于A类的内部类 , A类就属于B类的外部类

    ![image-20210402191054318](%E5%A4%9A%E6%80%81&%E5%B8%B8%E7%94%A8API.assets/image-20210402191054318.png)

- ##### 什么时候使用内部类 ?

  - 多个事物之间有包含关系, 可以使用内部类

- ##### 内部类分类 ? 

  - 成员内部类
  - 局部内部类
  - 匿名内部类

### 2.2 成员内部类

- 定义的位置 : 类中方法外

- 创建成员内部类对象格式 : 外部类名.内部类名  对象名 = new 外部类名().new 内部类名(参数);

  ```java
  package com.fan.innerclass_demo.member_innerclass;
  
  // 外部类
  public class Person {
      // 成员内部类
      public class Heart {
          // 频率变量
          private int rate;
          // 跳动方法
          public void beats() {
              System.out.println("咚咚咚!");
          }
      }
  }
  
  class Test {
      public static void main(String[] args) {
          // 创建内部类对象
          Person.Heart heart = new Person().new Heart();
          // 调用内部类中的方法
          heart.beats();
      }
  }
  ```

- ##### 成员内部类访问外部类的成员

  - 在内部类中有代表外部类对象的格式 : 外部类名的.this  , 私有的也可以访问
  - 外部类要想访问内部类成员 , 需要创建内部类对象

  ```java
  package com.fan.innerclass_demo.member_innerclass;
  
  public class Person {
      private String name = "张三";
      private int num = 10;
      
      // 成员内部类
      public class Heart {
          int num = 100;
          // 频率
          private int rate;
          // 跳动
          public void beats() {
              System.out.println("咚咚咚!");
          }
          // 调用外部类的成员
          public void show(){
              int num = 1000;
              System.out.println(Person.this.name);
              System.out.println(num);// 1000 就近原则
              System.out.println(this.num);// 100
              System.out.println(Person.this.num);// 10
  
          }
      }
  }
  
  class Test {
      public static void main(String[] args) {
          Person.Heart heart = new Person().new Heart();
          heart.beats();
          heart.show();
      }
  }
  ```

### 2.3 匿名内部类

- 匿名内部类 : 没有名字的类 ,  一次性产品
- 使用场景 : 直接调用方法 , 作为方法的传参 , 返回值类型 
- 好处 : 简化代码 , 快速实现接口或者抽象的抽象方法
- 格式 : 
  - new 类名/接口名(){  重写抽象方法 }   注意 : 此处创建的是子类对象!!!
- 使用方式 :
  - 直接调用方法
  - 作为方法的参数传递
  - 作为方法的返回值类型

```java
//接口
interface Flyable {
    void fly();
}
```

```java
// 直接调用方法
Flyable f1 = new Flyable() {
    @Override
    public void fly() {
        System.out.println("不知道什么在飞.....");
    }
};
f1.fly();
```

```java
// 作为方法的参数传递
showFlyable(
    new Flyable() {
        @Override
        public void fly() {
            System.out.println("不知道什么在飞3333");
        }
    }
);

public static void showFlyable(Flyable flyable) {
	flyable.fly();
}
```

```java
// 作为方法的返回值类型
public static Flyable getFlyable() {
	return new Flyable() {
        @Override
        public void fly() {
            System.out.println("3333333333333");
        }
    };
}
```

```java
package com.fan.innerclass_demo.anonymous_innerclass;

/*
    1 如果方法的参数是一个类的话 , 调用此方法需要传入此类的对象或者此类的子类对象
    2 如果方法的返回值类型是一个类的话 , 需要返回此类的对象 , 或者此类的子类对象

     3 如果方法的参数是一个接口的话 , 调用此方法需要传入此接口的实现类对象
     4 如果方法的返回值类型是一个接口的话 , 需要返回此接口的实现类对象


     匿名内部类 : 代表的就是子类对象!!!
        new 类名/接口名(){
            重写抽象类或者接口中的抽象方法
        };

     使用方向 :
        1 调用方法
        2 作为方法参数传递
        3 作为方法的返回值
 */
public interface Swim {
    public abstract void swimming();
}

class Test {
    public static void main(String[] args) {
//        // 子类对象!!!
        //  1 调用方法
//       new Swim() {
//            @Override
//            public void swimming() {
//                System.out.println("匿名内部类 , 重写了接口中的抽象方法...");
//            }
//        }.swimming();


//        //   2 作为方法参数传递
//        useSwim(new Swim() {
//            @Override
//            public void swimming() {
//                System.out.println("匿名内部类 , 重写了接口中的抽象方法...");
//            }
//        });

//        // 3 作为方法的返回值
//        Swim s = getSwim();
//        s.swimming();
    }


    public static Swim getSwim() {
        return new Swim() {
            @Override
            public void swimming() {
                System.out.println("匿名内部类 , 重写了接口中的抽象方法...");
            }
        };
    }


    /*
        Swim swim = new Swim() {
            @Override
            public void swimming() {
                System.out.println("匿名内部类 , 重写了接口中的抽象方法...");
            }
        };
     */
    public static void useSwim(Swim swim) {
        swim.swimming();
    }
}
```

## 3 API 

### 3.1 Object类

- 概述 : 类Object是类层次结构的根，每个类都把Object作为超类。 所有对象（包括数组）都实现了这个类的方法
- 方法 : public String toString()
  - 如果一个类没有重写toString方法 , 那么打印此类的对象 , 打印的是此对象的地址值
  - 如果一个类重写了toString方法 , 那么打印此类的对象 , 打印的是此对象的属性值
  - 好处 : 把对象转成字符串 , 快速查看一个对象的属性值
  - 执行时机 : 打印对象时会默认调用toString方法
- 方法 : public boolean equals()
  - 如果一个类没有重写equals方法 , 那么比较此类的对象 . 比较的是地址值
  - 如果一个类重写了equals方法 . 那么比较此类的对象 , 比较的是属性值是否相等
  - 好处 : 可以比较对象的内容

### 3.2 Objects类

- Objects是JDK1.7新增的一个对象工具类，里面都是静态方法可以用来操作对象。比如对象的比较，计算对象的hash值，判断对手是否为空....比如里面的equals方法，可以避免空指针异常

```java
public static boolean equals(Object a, Object b):判断两个对象是否相等
    
public static boolean equals(Object a, Object b) {
    return (a == b) || (a != null && a.equals(b));
}
a.equals(b) :如果a是null值，肯定会空指针
Objects.equals(a,b);：如果a是null，不会导致空指针异常
```

### 3.2  Date类

- ##### 概述 : java.util.Date 表示特定的瞬间，精确到毫秒

- ##### 构造方法 : 

  - public Date(): 创建的对象，表示的是当前计算机系统的时间
  - public Date(long time): Date对象 = 1970/1/1 0:0:0 + time毫秒值

- ##### 成员方法 : 

  - public long getTime(): 返回毫秒值 = 当前Date代表的时间 - 1970/1/1 0:0:0
  - public void setTime(long t): Date对象 = 1970/1/1 0:0:0 + time毫秒值

```java
package com.fan.api_demo.date_demo;

import java.util.Date;

/*
    Date类 : 代表的是一个瞬间 , 精确到毫秒

    构造方法 :
        public Date() : 代表的是当前系统时间
        public Date(long date) : Date对象 = 1970/1/1 0:0:0 + long类型的毫秒值

    成员方法 :
        public void setTime(long date) : Date对象 = 1970/1/1 0:0:0 + long类型的毫秒值
        public long getTime() : 返回的是毫秒值 = Date代表的时间 - 1970/1/1 0:0:0
 */
public class DateDemo {
    public static void main(String[] args) {
        //  public Date() : 代表的是当前系统时间
//        Date d = new Date();
//        System.out.println(d);

        //  public Date(long date) : Date对象 = 1970/1/1 0:0:0 + long类型的毫秒值
//        Date d2 = new Date(1000L * 60 * 60 * 24); // 1970/1/1 0:0:0 + 一天的毫秒值
//        System.out.println(d2);

        Date d = new Date();
        // public void setTime(long date) : Date对象 = 1970/1/1 0:0:0 + long类型的毫秒值
        // d.setTime(1000L * 60 * 60 * 24);
        System.out.println(d);

        // public long getTime() : 返回的是毫秒值 = Date代表的时间 - 1970/1/1 0:0:0
        // System.out.println(d.getTime());

    }
}
```

### 3.3 DateFormat类

- 概述 : 主要用于操作日期格式的一个类

- 作用 :

  - 格式化 : Date --> String
  - 解析 : String --> Date

- 构造方法 : 

  - SimpleDateFormat(String  pattern)    给定日期模板创建日期格式化对象

    ![image-20210402212755339](%E5%A4%9A%E6%80%81&%E5%B8%B8%E7%94%A8API.assets/image-20210402212755339.png)

- 成员方法 : 

  - public String format ( Date d )：格式化，将日期对象格式化为字符串
  - public Date parse ( String s )：解析，将字符串解析为日期对象

```java
package com.fan.api_demo.dateformat_demo;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

/*
    SimpleDateFormat类 :
        "2021年04月03日 16:48:10"  ---> Date
        Date(Sat Apr 03 16:41:38 CST 2021) --> "2021年04月03日 16:48:10"

        作用 :
            格式化 : Date --> String
            解析   : String --> Date

        构造方法 :
            public SimpleDateFormat(String pattern) : pattern : 字符串类型的日期模板

        成员方法 ;
            public final String format(Date date) : 接收一个Date对象返回指定模式的字符串
            public Date parse(String source) : 接收一个字符串  , 返回一个Date对象



        1 获取当前的日期对象，使用格式：yyyy-MM-dd HH:mm:ss 来表示，例如：2020-10-31 17:00:00【格式化】

        2 将字符串的 2020年10月31日  17:00:00，转换为日期Date对象。【解析】

 */
public class SimpleDateFormatDemo {
    public static void main(String[] args) throws ParseException {
        // 解析   : String --> Date
        String strDate = "2020年10月31日 17:00:00";
        // 注意 : 解析时 , SimpleDateFormat的参数(日期模板) , 必须和要解析字符串的模式匹配
        SimpleDateFormat sdf = new SimpleDateFormat("yyyy年MM月dd日 HH:mm:ss");
        //  public Date parse(String source) : 接收一个字符串  , 返回一个Date对象
        Date date = sdf.parse(strDate);
        System.out.println(date); // Sat Oct 31 17:00:00 CST 2020
    }

    private static void method() {
        // 格式化 : Date --> String
        // 获取当前系统时间
        Date date = new Date();
        // System.out.println(date);// Sat Apr 03 16:53:35 CST 2021

        // public SimpleDateFormat(String pattern) : pattern : 字符串类型的日期模板
        SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");

        // public final String format(Date date) : 接收一个Date对象返回指定模式的字符串
        String strDate = sdf.format(date);

        System.out.println(strDate);
    }
}
```

### 3.4 Calendar类

- ##### 概述 : 

  - java.util.Calendar类表示日历，内含有各种时间信息，以及获取，计算时间的方法。
  - Calendar本身是一个抽象类，可以通过Calendar提供的静态方法getInstance日历子类对象

- ##### Calendar常用方法 : 

  | **方法名**                            | **说明**                                |
  | ------------------------------------- | --------------------------------------- |
  | public static Calendar getInstance()  | 获取一个它的子类GregorianCalendar对象。 |
  | public int get(int field)             | 获取指定日历字段的时间值。              |
  | public void set(int field,int value)  | 设置指定字段的时间值                    |
  | public void add(int field,int amount) | 为某个字段增加/减少指定的值             |

- get,set,add方法参数中都有一个整数类型的参数field，field用来区分当前要获取或者操作的日期时间数据的。field数据的值使用Calender中定义的整数常量来表示

  - Calendar.YEAR : 年 
  - Calendar.MONTH ：月  
  - Calendar.DAY_OF_MONTH：月中的周
  - Calendar.HOUR：小时 
  - Calendar.MINUTE：分钟，
  - Calendar.SECOND：秒
  - Calendar.DAY_OF_WEEK：星期

- ##### 星期对应的关系

  ​     字段时间值  :   1           		2            			3          		  ...        7
  ​     真正的时间  : SUNDAY        MONDAY      TUESDAY    	 ...        SATURDAY

- ##### 月份对应的关系

  ​    字段时间值     :       0           				1             ....        	11
  ​    正真的时间     :     JANUARY     FEBRUARY       ....      	DECEMBER

```java
package com.fan.api_demo.calendar_demo;
/*

    月份对应的关系
    字段时间值     :       0           1           ....        11
    正真的时间     :     JANUARY     FEBRUARY      ....      DECEMBER


 */
import java.util.Calendar;
/*
    成员方法 :
        public int get(int field)	获取指定日历字段的时间值。
        public void set(int field,int value)	设置指定字段的时间值
        public void add(int field,int amount)	为某个字段增加/减少指定的值
 */
public class CalendarDemo {
    public static void main(String[] args) {
        // 获取Calendar对象 , rightNow对象
        Calendar rightNow = Calendar.getInstance();

        // public void set(int field , int value)	设置指定字段的时间值
        // rightNow.set(Calendar.YEAR , 2024);
        // rightNow.set(Calendar.MONTH , 5);
        // rightNow.set(Calendar.DAY_OF_MONTH, 10);

        // public void add(int field,int amount)	为某个字段增加/减少指定的值
        // rightNow.add(Calendar.DAY_OF_MONTH , -3);

        //  public int get(int field)	 : 获取指定日历字段的时间值。
        int year = rightNow.get(Calendar.YEAR);
        int month = rightNow.get(Calendar.MONTH);
        int day = rightNow.get(Calendar.DAY_OF_MONTH);


        System.out.println(year); // 2021
        System.out.println(month);// 3
        System.out.println(day);  // 3
    }
}
```

```java
package com.fan.api_demo.calendar_demo;

import java.util.Calendar;

/*
    1 写代码使用get方法，将年，月，日，时，分，秒，周获取出来
    特别注意获取月份，和星期有以下特点：
    直接获取的月份数据是从0开始的， 0表示1月，.....11表示12月
    周的获取，从周日开始计算，1就是周日，2就是周一 ......7就是周六


星期对应的关系
     字段时间值  :   1           2            3          ...        7
     真正的时间  : SUNDAY        MONDAY      TUESDAY     ...        SATURDAY

月份对应的关系
    字段时间值     :       0           1           ....        11
    正真的时间     :     JANUARY     FEBRUARY      ....      DECEMBER

 */
public class Test1 {
    public static void main(String[] args) {
        // 拿到当前时间
        Calendar now = Calendar.getInstance();

        System.out.println(now.get(Calendar.YEAR));
        System.out.println(now.get(Calendar.MONTH));
        System.out.println(now.get(Calendar.DAY_OF_MONTH));
        System.out.println(now.get(Calendar.HOUR));
        System.out.println(now.get(Calendar.MINUTE));
        System.out.println(now.get(Calendar.SECOND));

        int week = now.get(Calendar.DAY_OF_WEEK);// 7
        System.out.println(getWeek(week));// 字段值
    }

    public static String getWeek(int field){
        String[] str = { "" , "SUNDAY" , "MONDAY" , "TUESDAY" , "WEDNESDAY" , "THURSDAY" , "FRIDAY" , "SATURDAY"};
        return str[field];
    }
}
```

```java
package com.fan.api_demo.calendar_demo;

import java.util.Calendar;

/*
    2 写代码实现，获取2022年10月1日是星期几？
    参考思路：
    直接获取日历对象，得到的是当前系统的日历时间信息。
    获取日历对象后，要重新设置日期
    再获取星期数据

 */
public class Test2 {
    public static void main(String[] args) {
        Calendar cal = Calendar.getInstance();
        cal.set(Calendar.YEAR, 2022);
        cal.set(Calendar.MONTH, 9);
        cal.set(Calendar.DAY_OF_MONTH, 1);

        int week = cal.get(Calendar.DAY_OF_WEEK);
        System.out.println(getWeek(week));// 字段值
    }

    public static String getWeek(int field){
        String[] str = { "" , "SUNDAY" , "MONDAY" , "TUESDAY" , "WEDNESDAY" , "THURSDAY" , "FRIDAY" , "SATURDAY"};
        return str[field];
    }
}
```

```java
package com.fan.api_demo.calendar_demo;

import java.util.Calendar;

/*
    3 计算10000天之后的年月日

    参考思路：
    先获取当前日历对象
    再调用add方法，指定DATE或者DAY_OF_MONTH，添加10000天
    再获取日历的年，月，日

 */
public class Test3 {
    public static void main(String[] args) {
        Calendar cal = Calendar.getInstance();
        cal.add(Calendar.DAY_OF_MONTH, 10000);


        System.out.println(cal.get(Calendar.YEAR));
        System.out.println(cal.get(Calendar.MONTH));
        System.out.println(cal.get(Calendar.DAY_OF_MONTH));
    }
}
```

### 3.5 Math类

- 概述 : Math包含执行基本数字运算的方法，如基本指数，对数，平方根和三角函数。所提供的都是静态方法，可以直接调用

- 常用方法 : 

  | public static int abs(int a)                 | 获取参数a的绝对值： |
  | -------------------------------------------- | ------------------- |
  | public static double ceil(double a)          | 向上取整            |
  | public static double floor(double a)         | 向下取整            |
  | public static double pow(double a, double b) | 获取a的b次幂        |
  | public static long round(double a)           | 四舍五入取整        |

### 3.6 System类 

- System类包含几个有用的类字段和方法。 它不能被实例化

- 常用方法 : 

  | **方法名**                             | **说明**                                     |
  | -------------------------------------- | -------------------------------------------- |
  | public static void exit(int status)    | 终止当前运行的 Java 虚拟机，非零表示异常终止 |
  | public static long currentTimeMillis() | 返回当前时间(以毫秒为单位)                   |

## 4 BigInteger类

### 4.1 概述

- 概述 : java.math.BigInteger类是一个引用数据类型 , 可以用于计算一些大的整数 , 当超出基本数据类型数据范围的整数运算时就可以使用BigInteger了。

### 4.2 构造方法

- 构造方法 : 可以将整数的字符串 . 转成BigInteger类型的对象

### 4.3 成员方法

- | **方法声明**                                      | **描述**                           |
  | ------------------------------------------------- | ---------------------------------- |
  | public BigInteger **add** (BigInteger value)      | 超大整数加法运算                   |
  | public BigInteger **subtract** (BigInteger value) | 超大整数减法运算                   |
  | public BigInteger **multiply** (BigInteger value) | 超大整数乘法运算                   |
  | public BigInteger **divide** (BigInteger value)   | 超大整数除法运算，除不尽取整数部分 |

```java
package com.fan.api_demo.biginteger_demo;

import java.math.BigInteger;

/*
    构造方法 :
        BigInteger(String value)	可以将整数的字符串，转换为BigInteger对象
    成员方法 :
        public BigInteger add (BigInteger value)	    超大整数加法运算
        public BigInteger subtract (BigInteger value)	超大整数减法运算
        public BigInteger multiply (BigInteger value)	超大整数乘法运算
        public BigInteger divide (BigInteger value)	超大整数除法运算，除不尽取整数部分

 */
public class BigIntegerDemo {
    public static void main(String[] args) {
        // 获取大整数对象
        BigInteger bigInteger1 = new BigInteger("2147483648");
        // 获取大整数对象
        BigInteger bigInteger2 = new BigInteger("2");
        // public BigInteger add (BigInteger value)	    超大整数加法运算
        BigInteger add = bigInteger1.add(bigInteger2);
        System.out.println(add);

        System.out.println("=============");

        // public BigInteger subtract (BigInteger value)	超大整数减法运算
        BigInteger subtract = bigInteger1.subtract(bigInteger2);
        System.out.println(subtract);

        System.out.println("=============");

        // public BigInteger multiply (BigInteger value)	超大整数乘法运算
        BigInteger multiply = bigInteger1.multiply(bigInteger2);
        System.out.println(multiply);

        System.out.println("=============");
        // public BigInteger divide (BigInteger value)	超大整数除法运算，除不尽取整数部分
        BigInteger divide = bigInteger1.divide(bigInteger2);
        System.out.println(divide);
    }
}
```

## 5 BigDecimal类

### 5.1 概述

- 概述 : java.math.BigDecimal可以对大浮点数进行运算，保证运算的准确性。float，double 他们在存储及运算的时候，会导致数据精度的丢失。如果要保证运算的准确性，就需要使用BigDecimal。

### 5.2 构造方法

- 构造方法 : 
  - public BigDecimal(String val) : 将 BigDecimal 的字符串表示形式转换为 BigDecimal

### 5.3 成员方法

- 成员方法 : 

  - | **方法声明**                                                 | **描述**                                                     |
    | ------------------------------------------------------------ | ------------------------------------------------------------ |
    | public BigDecimal **add**(BigDecimal value)                  | 加法运算                                                     |
    | public BigDecimal **subtract**(BigDecimal value)             | 减法运算                                                     |
    | public BigDecimal **multiply**(BigDecimal value)             | 乘法运算                                                     |
    | public BigDecimal **divide**(BigDecimal value)               | 除法运算(除不尽会有异常)                                     |
    | public BigDecimal divide(BigDecimal divisor, int roundingMode) | 除法运算(除不尽，使用该方法)参数说明：scale 精确位数，roundingMode取舍模式         BigDecimal.ROUND_HALF_UP 四舍五入      BigDecimal.ROUND_FLOOR 去尾法      BigDecimal.ROUND_UP 进一法 |

```java
package com.fan.api_demo.bigdecimal_demo;

import java.math.BigDecimal;

/*
    构造方法 :
        public BigDecimal(String val)	将 BigDecimal 的字符串表示形式转换为 BigDecimal
    成员方法 :
        public BigDecimal add(BigDecimal value)	加法运算
        public BigDecimal subtract(BigDecimal value)	减法运算
        public BigDecimal multiply(BigDecimal value)	乘法运算
        public BigDecimal divide(BigDecimal value)	除法运算(除不尽会有异常)
        public BigDecimal divide(BigDecimal value, int scale, RoundingMode roundingMode)	除法运算(除不尽，使用该方法)
        参数说明：
        scale 精确位数，
        roundingMode取舍模式
                   BigDecimal.ROUND_HALF_UP 四舍五入
                   BigDecimal.ROUND_FLOOR 去尾法
                   BigDecimal.ROUND_UP  进一法
 */
public class BigDecimalDemo {
    public static void main(String[] args) {
        BigDecimal bigDecimal1 = new BigDecimal("0.1");
        BigDecimal bigDecimal2 = new BigDecimal("0.2");

        // public BigDecimal add(BigDecimal value)	加法运算
        BigDecimal add = bigDecimal1.add(bigDecimal2);
        System.out.println(add);

        System.out.println("=================");

        // public BigDecimal subtract(BigDecimal value)	减法运算
        BigDecimal subtract = bigDecimal1.subtract(bigDecimal2);
        System.out.println(subtract);

        System.out.println("=================");

        // public BigDecimal multiply(BigDecimal value)	乘法运算
        BigDecimal multiply = bigDecimal1.multiply(bigDecimal2);
        System.out.println(multiply);

        System.out.println("=================");

        // public BigDecimal divide(BigDecimal value)	除法运算(除不尽会有异常)
        // BigDecimal divide = bigDecimal1.divide(bigDecimal2);
        // System.out.println(divide);

        /*
            public BigDecimal divide(BigDecimal divisor, int roundingMode)	除法运算(除不尽，使用该方法)
            参数说明：
            scale 精确位数，
            roundingMode : 取舍模式
                       BigDecimal.ROUND_HALF_UP 四舍五入
                       BigDecimal.ROUND_FLOOR 去尾法
                       BigDecimal.ROUND_UP  进一法
        */

        // BigDecimal divide = bigDecimal1.divide(bigDecimal2, 3, BigDecimal.ROUND_HALF_UP);
        // BigDecimal divide = bigDecimal1.divide(bigDecimal2, 3, BigDecimal.ROUND_FLOOR);
        // BigDecimal divide = bigDecimal1.divide(bigDecimal2, 3, BigDecimal.ROUND_UP);
        // System.out.println(divide);

    }
}
```

## 6 Arrays类

### 6.1 概述

- 概述 : java.util.Arrays是数组的工具类，里面有很多静态的方法用来对数组进行操作（如排序和搜索），还包含一个静态工厂，可以将数组转换为List集合（后面会讲到集合知识

### 6.2 构造方法

- 构造方法 : private Arrays(){}

- | **public static void sort(int[] a)**       | **按照数字顺序排列指定的数组**         |
  | ------------------------------------------ | -------------------------------------- |
  | **public static String toString(int[] a)** | **返回指定数组的内容的字符串表示形式** |

```java
package com.fan.api_demo.arrays_demo;

import java.util.Arrays;
import java.util.Random;

/*
    1 随机生成10个[0,100]之间的整数，放到整数数组中，按照从小到大的顺序排序并打印元素内容。
 */
public class ArraysDemo {
    public static void main(String[] args) {
        // 创建数组
        int[] arr = new int[10];

        // 创建随机数对象
        Random r = new Random();

        // 采用随机数给数组的每一个元素赋值
        for (int i = 0; i < arr.length; i++) {
            arr[i] = r.nextInt(101);
        }

        // 对数组进行排序
        Arrays.sort(arr);

        // 把数组转成指定格式的字符串
        System.out.println(Arrays.toString(arr));
    }
}

```

## 7 包装类

### 7.1 概述

- ##### 概述 :

  - Java中基本数据类型对应的引用数据类型

### 7.2 包装类的作用

- ##### 包装类的作用 : 

  - 基本数据类型 , 没有变量 , 没有方法 , 包装类就是让基本数据类型拥有变量和属性 , 实现对象化交互
  - 基本数据类型和字符串之间的转换

### 7.3 基本数据类型和包装类对应

- ##### 基本数据类型和包装类的对应关系

  | **基本数据类型** | **包装类型** |
  | ---------------- | ------------ |
  | byte             | Byte         |
  | short            | Short        |
  | int              | Integer      |
  | long             | Long         |
  | float            | Float        |
  | double           | Double       |
  | char             | Character    |
  | boolean          | Boolean      |

### 7.4 自动装箱和自动拆箱

- 自动转型和自动拆箱

  - 自动装箱和拆箱是JDK1.5开始的
  - 自动装箱 : 基本数据类型自动转成对应的包装类类型
  - 自动拆箱 : 包装类类型自动转成对应的基本数据类型

  ```java
  Integer i1 = 10;
  int i2 = i1;
  ```

### 7.5 基本数据类型和字符串之间的转换

- 使用包装类, 对基本数据类型和字符串之间的转换

  - 在开发过程中数据在不同平台之间传输时都以字符串形式存在的，有些数据表示的是数值含义，如果要用于计算我们就需要将其转换基本数据类型.

  - 基本数据类型--> String

    - 直接在数值后加一个空字符串
    - 通过String类静态方法valueOf()

  - String --> 基本数据类型

    - | public static byte parseByte(String s)：将字符串参数转换为对应的byte基本类型。 |
      | ------------------------------------------------------------ |
      | public static short parseShort(String s)：将字符串参数转换为对应的short基本类型。 |
      | public static int parseInt(String s)：将字符串参数转换为对应的int基本类型。 |
      | public static long parseLong(String s)：将字符串参数转换为对应的long基本类型。 |
      | public static float parseFloat(String s)：将字符串参数转换为对应的float基本类型。 |
      | public static double parseDouble(String s)：将字符串参数转换为对应的double基本类型。 |
      | public static boolean parseBoolean(String s)：将字符串参数转换为对应的boolean基本类型。 |

- 注意事项 : 

  - 包装类对象的初始值为null（是一个对象）
  - Java中除了float和double的其他基本数据类型，都有常量池
    - 整数类型：[-128,127]值在常量池
    - 字符类型：[0,127]对应的字符在常量池
    - 布尔类型：true，false在常量池

  - 在常量池中的数据 , 会进行共享使用 , 不在常量池中的数据会创建一个新的对象

## 8 String类的常用方法

### 8.1 常用方法

![image-20210404212729376](%E5%A4%9A%E6%80%81&%E5%B8%B8%E7%94%A8API.assets/image-20210404212729376.png)

```java
package com.fan.api_demo.string_demo;

/*
    已知字符串，完成需求
    String str = "I Love Java, I Love Heima";
    判断是否存在  “Java”
    判断是否以Heima字符串结尾
    判断是否以Java开头
    判断 Java在字符串中的第一次出现位置
    判断  itcast 所在的位置

 */
public class Test {
    public static void main(String[] args) {
        String str = "I Love Java, I Love Heima";

        // 判断是否存在  “Java”
        System.out.println(str.contains("Java"));// true

        // 判断是否以Heima字符串结尾
        System.out.println(str.endsWith("Heima"));// true

        // 判断是否以Java开头
        System.out.println(str.startsWith("Java"));// false

        // 判断 Java在字符串中的第一次出现位置
        System.out.println(str.indexOf("Java"));// 7

        // 判断  itcast 所在的位置
        System.out.println(str.indexOf("itcast"));// -1
    }
}

```

```java
package com.fan.api_demo.string_demo;

/*
    已知字符串，完成右侧需求
    String str = "I Love Java, I Love Heima";
    需求：
    1.将所有 Love 替换为 Like ,打印替换后的新字符串
    2.截取字符串 "I Love Heima"
    3.截取字符串 "Java"

 */
public class Test2 {
    public static void main(String[] args) {
        String str = "I Love Java, I Love Heima";

        // 1.将所有 Love 替换为 Like ,打印替换后的新字符串
        System.out.println(str.replace("Love", "Like"));
        // 2.截取字符串 "I Love Heima"
        System.out.println(str.substring(13));
        // 3.截取字符串 "Java"
        System.out.println(str.substring(7 , 11));

    }
}

```

```java
package com.fan.api_demo.string_demo;

/*
    已知字符串，完成右侧需求
    String str = "I Love Java, I Love Heima";
    需求：
    1 计算字符 a 出现的次数（要求使用toCharArray）
    2 计算字符 a 出现的次数（要求使用charAt）
    3 将字符串中所有英文字母变成小写
    4 将字符串中所有英文字母变成大写
 */
public class Test3 {
    public static void main(String[] args) {
        String str = "I Love Java, I Love Heima";

//        1 计算字符 a 出现的次数（要求使用toCharArray）
        int count1 = 0;
        char[] chars = str.toCharArray();
        for (int i = 0; i < chars.length; i++) {
            if (chars[i] == 'a') {
                count1++;
            }
        }
        System.out.println("字符a出现了" + count1 + "次");


//        2 计算字符 a 出现的次数（要求使用charAt）
        int count2 = 0;
        for (int i = 0; i < str.length(); i++) {
            char ch = str.charAt(i);
            if(ch == 'a'){
                count2++;
            }
        }
        System.out.println("字符a出现了" + count2 + "次");

//        3 将字符串中所有英文字母变成小写
        String s1 = str.toLowerCase();
        System.out.println(s1);

//        4 将字符串中所有英文字母变成大写
        String s2 = str.toUpperCase();
        System.out.println(s2);
    }
}
```

## 9 正则表达式

### 9.1 概述 : 

- 正则表达式通常用来校验，检查字符串是否符合规则的

### 9.2 体验正则表达式

```java
package com.fan.regex_demo;

import java.util.Scanner;

/*
    设计程序让用户输入一个QQ号码，验证QQ号的合法性：
    1. QQ号码必须是5--15位长度
    2. 而且首位不能为0
    3. 而且必须全部是数字

 */
public class Test1 {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);

        System.out.println("请输入您的qq号码:");
        String qq = sc.nextLine();

        System.out.println(checkQQ2(qq));

    }

    private static boolean checkQQ(String qq) {
//        1. QQ号码必须是5--15位长度
        if (qq.length() < 5 || qq.length() > 15) {
            return false;
        }
//       2 . 而且首位不能为0
        if (qq.charAt(0) == '0') {
            return false;
        }

//        2. 而且必须全部是数字
        for (int i = 0; i < qq.length(); i++) {
            char ch = qq.charAt(i);
            if (ch < '0' || ch > '9') {
                return false;
            }
        }

        return true;
    }

    // 正则表达式改进
    private static boolean checkQQ2(String qq) {
        return qq.matches("[1-9][0-9]{4,14}");
    }
}

```

### 6.3 正则表达式的语法

- ##### boolean  matches（正则表达式） :如果匹配正则表达式就返回true，否则返回false

  - boolean  matches（正则表达式） :如果匹配正则表达式就返回true，否则返回false

- ##### 字符类

  -  [abc] ：代表a或者b，或者c字符中的一个。
  -  [^abc]：代表除a,b,c以外的任何字符。
  -  [a-z] ：代表a-z的所有小写字符中的一个。
  -  [A-Z] ：代表A-Z的所有大写字符中的一个。
  -  [0-9] ：代表0-9之间的某一个数字字符。
  -  [a-zA-Z0-9]：代表a-z或者A-Z或者0-9之间的任意一个字符。
  -  [a-dm-p]：a 到 d 或 m 到 p之间的任意一个字符

```java
package com.fan.regex_demo;
/*
    字符类 : 方括号被用于指定字符
    [abc] ：代表a或者b，或者c字符中的一个。
    [^abc]：代表除a,b,c以外的任何字符。
    [a-z] ：代表a-z的所有小写字符中的一个。
    [A-Z] ：代表A-Z的所有大写字符中的一个。
    [0-9] ：代表0-9之间的某一个数字字符。
    [a-zA-Z0-9]：代表a-z或者A-Z或者0-9之间的任意一个字符。
    [a-dm-p]：a 到 d 或 m 到 p之间的任意一个字符

    需求 :
    1 验证str是否以h开头，以d结尾，中间是a,e,i,o,u中某个字符
    2 验证str是否以h开头，以d结尾，中间不是a,e,i,o,u中的某个字符
    3 验证str是否a-z的任何一个小写字符开头，后跟ad
    4 验证str是否以a-d或者m-p之间某个字符开头，后跟ad
    注意: boolean  matches（正则表达式） :如果匹配正则表达式就返回true，否则返回false
 */
public class RegexDemo {
    public static void main(String[] args) {
//        1 验证str是否以h开头，以d结尾，中间是a,e,i,o,u中某个字符
        System.out.println("had".matches("h[aeiou]d"));

//        2 验证str是否以h开头，以d结尾，中间不是a,e,i,o,u中的某个字符
        System.out.println("hwd".matches("h[^aeiou]d"));

//        3 验证str是否a-z的任何一个小写字符开头，后跟ad
        System.out.println("aad".matches("[a-z]ad"));

//        4 验证str是否以a-d或者m-p之间某个字符开头，后跟ad
        System.out.println("bad".matches("[a-dm-p]ad"));

    }
}
```

- ##### 逻辑运算符

  - && ：并且
  - |    ：或者

```java
package com.fan.regex_demo;
/*
    逻辑运算符 :
        1 && : 并且
        2 |  : 或者

    需求 :
        1 要求字符串是除a、e、i、o、u外的其它小写字符开头，后跟ad
        2 要求字符串是aeiou中的某个字符开头，后跟ad
 */
public class RegexDemo2 {
    public static void main(String[] args) {
//        1 要求字符串是除a、e、i、o、u外的其它小写字符开头，后跟ad
        System.out.println("vad".matches("[a-z&&[^aeiou]]ad"));
//        2 要求字符串是aeiou中的某个字符开头，后跟ad
        System.out.println("aad".matches("[a|e|i|o|u]ad"));
    }
}
```

- ##### 预定义字符类

  - "."  ： 匹配任何字符。
  - "\d"：任何数字[0-9]的简写；
  - "\D"：任何非数字[^0-9]的简写；
  - "\s" ： 空白字符：[ \t\n\x0B\f\r] 的简写
  - "\S" ： 非空白字符：[^\s] 的简写
  - "\w" ：单词字符：[a-zA-Z_0-9]的简写
  - "\W"：非单词字符：[^\w]

```java
package com.fan.regex_demo;
/*
    预定义字符 : 简化字符类的书写

    "."  ：匹配任何字符。
    "\d" ：任何数字[0-9]的简写
    "\D" ：任何非数字[^0-9]的简写
    "\s" ：空白字符：[\t\n\x0B\f\r] 的简写
    "\S" ：非空白字符：[^\s] 的简写
    "\w" ：单词字符：[a-zA-Z_0-9]的简写
    "\W" ：非单词字符：[^\w]

    需求 :
       1 验证str是否3位数字
       2 验证手机号：1开头，第二位：3/5/8，剩下9位都是0-9的数字
       3 验证字符串是否以h开头，以d结尾，中间是任何字符

 */
public class RegexDemo3 {
    public static void main(String[] args) {
//        1 验证str是否3位数字
        System.out.println("123".matches("\\d\\d\\d"));

//        2 验证手机号：1开头，第二位：3/5/8，剩下9位都是0-9的数字 ）
        System.out.println("15188888888".matches("1[358]\\d\\d\\d\\d\\d\\d\\d\\d\\d"));

//        3 验证字符串是否以h开头，以d结尾，中间是任何字符
        System.out.println("had".matches("h.d"));
    }
}
```

- ##### 数量词

  - X? : 0次或1次
  - X* : 0次到多次
  - X+ : 1次或多次
  - X{n} : 恰好n次
  - X{n,} : 至少n次
  - X{n,m}: n到m次(n和m都是包含的)

```java
package com.fan.regex_demo;

/*
    数量词 :
        - X?    : 0次或1次
        - X*    : 0次到多次
        - X+    : 1次或多次
        - X{n}  : 恰好n次
        - X{n,} : 至少n次
        - X{n,m}: n到m次(n和m都是包含的)

    需求 :
      1 验证str是否3位数字
      2 验证str是否是多位(大于等于1次)数字
      3 验证str是否是手机号 ( 1开头，第二位：3/5/8，剩下9位都是0-9的数字)
      4 验证qq号码：1).5--15位；2).全部是数字;3).第一位不是0

 */
public class RegexDemo4 {
    public static void main(String[] args) {
//        1 验证str是否3位数字
        System.out.println("123".matches("\\d{3}"));

//        2 验证str是否是多位(大于等于1次)数字
        System.out.println("123456".matches("\\d+"));

//        3 验证str是否是手机号 ( 1开头，第二位：3/5/8，剩下9位都是0-9的数字)
        System.out.println("15188888888".matches("1[358]\\d{9}"));

//        4 验证qq号码：1).5--15位；2).全部是数字;3).第一位不是0
        System.out.println("122103987".matches("[1-9]\\d{4,14}"));
    }
}
```

- ##### 分组括号 : 

  - 将要重复使用的正则用小括号括起来，当做一个小组看待

```java
package com.fan.regex_demo;
/*
    分组括号 : 将要重复使用的正则用小括号括起来，当做一个小组看待
    需求 :  window秘钥 , 分为5组，每组之间使用 - 隔开 , 每组由5位A-Z或者0-9的字符组成 , 最后一组没有 -
    举例 :
        xxxxx-xxxxx-xxxxx-xxxxx-xxxxx
        DG8FV-B9TKY-FRT9J-99899-XPQ4G
    分析：
        前四组其一  ：DG8FV-    正则：[A-Z0-9]{5}-
        最后一组    ：XPQ4G     正则：[A-Z0-9]{5}

    结果 : ([A-Z0-9]{5}-){4}[A-Z0-9]{5}

 */
public class RegexDemo5 {
    public static void main(String[] args) {
        System.out.println("DG8FV-B9TKY-FRT9J-99899-XPQ4G".matches("([A-Z0-9]{5}-){4}[A-Z0-9]{5}"));
    }
}
```

- ##### 字符串中常用含有正则表达式的方法

  - public String[] split ( String regex ) 可以将当前字符串中匹配regex正则表达式的符号作为"分隔符"来切割字符串。
  - public String replaceAll ( String regex , String newStr ) 可以将当前字符串中匹配regex正则表达式的字符串替换为newStr。

```java
package com.fan.regex_demo;

import java.util.Arrays;

/*

    1 字符串中常用含有正则表达式的方法
        public String[] split ( String regex ) 可以将当前字符串中匹配regex正则表达式的符号作为"分隔符"来切割字符串。
        public String replaceAll ( String regex , String newStr ) 可以将当前字符串中匹配regex正则表达式的字符串替换为newStr。

    需求:
        1 将以下字符串按照数字进行切割
        String str1 = "aa123bb234cc909dd";

        2 将下面字符串中的"数字"替换为"*“a
        String str2 = "我卡里有100000元，我告诉你卡的密码是123456，要保密哦";

 */
public class RegexDemo6 {
    public static void main(String[] args) {
        // 1 将以下字符串按照数字进行切割
        String str1 =  "aa123bb234cc909dd";
        String[] strs = str1.split("\\d+");
        System.out.println(Arrays.toString(strs));

        // 2 将下面字符串中的"数字"替换为"*“a
        String str2 = "我卡里有100000元，我告诉你卡的密码是123456，要保密哦";
        System.out.println(str2.replaceAll("\\d+" , "*"));
    }
}
```
