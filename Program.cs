using System;
using System.Drawing;
using System.Windows.Forms;

class Program
{
    [STAThread]
    static void Main()
    {
        Application.EnableVisualStyles();
        Application.Run(new Merlin());
    }
}

class Merlin : Form
{
    public Merlin()
    {
        this.Text = "merlin";
        this.Size = new Size(1000, 800);
        this.StartPosition = FormStartPosition.CenterScreen;

        Button btn = new Button
        {
            Text = "Click Me!",
            Size = new Size(150, 50),
            Location = new Point(125, 60),
            Font = new Font("Segoe UI", 12, FontStyle.Bold)
        };

        btn.Click += (sender, e) => 
        {
            MessageBox.Show("Hello World!", "Message");
        };

        this.Controls.Add(btn);
    }
}
