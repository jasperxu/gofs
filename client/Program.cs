using System;

namespace NetCoreHttp
{
    class Program
    {
        static void Main(string[] args)
        {
            var g = new GoFs { URL = "http://localhost:8080/" };

            if (g.Upload("微信截图_20170912143114.png", "/CS/2.png"))
                Console.WriteLine("Upload Success");
            else
                Console.WriteLine("Upload Error");

            if (g.Download("/CS/2.png", "aaa.png"))
                Console.WriteLine("Download Success");
            else
                Console.WriteLine("Download Error");

            if (g.Delete("/CS/2.png"))
                Console.WriteLine("Delete Success");
            else
                Console.WriteLine("Delete Error");

            Console.Read();
        }
    }
}
