import java.util.Scanner;

public class NameAge{
    public static void main(String[] args) {
        Scanner console = new Scanner(System.in);
        String name = console.nextLine();
        int age = console.nextInt();
        console.close();
        System.out.println("Name: " + name);
        System.out.println("Age: " + age);
    }
}