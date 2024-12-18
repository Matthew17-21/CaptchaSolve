# CaptchaSolve
Program to help streamline captcha solving by routing requests to third-party services like 2Captcha!

## TODO

- [ ] Initialize a solver for each site in the config
- [ ] Store captcha tokens until requested
- [ ] Scheduler to delete expired tokens?
- [ ] Each time a captcha is requested, and no tokens are available, start a goroutine 
- [x] Global solver vs instances
- [ ] If a API key in a given site is out of funds, do not keep requesting from the site
