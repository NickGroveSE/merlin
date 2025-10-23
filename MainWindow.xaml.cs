using System.Diagnostics;
using System.Windows;
using System.Windows.Controls;

namespace MerlinWpf
{
    public partial class MainWindow : Window
    {
        public MainWindow()
        {
            InitializeComponent();
        }


        private void Button_Click(object sender, RoutedEventArgs e)
        {
            MessageBox.Show("Hello World!", "Message");
        }

        private void OpenLink_Click(object sender, RoutedEventArgs e)
        {
            Button btn = (Button)sender;
            string url = btn.Tag.ToString();
            
            Process.Start(new ProcessStartInfo
            {
                FileName = url,
                UseShellExecute = true
            });
        }

    }
}