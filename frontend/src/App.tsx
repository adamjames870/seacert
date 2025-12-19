import { useState, useContext, useEffect } from 'react'
import { Routes, Route, Link as RouterLink, useNavigate } from 'react-router-dom'
import { 
  Typography, 
  Button, 
  AppBar, 
  Toolbar, 
  IconButton, 
  Menu, 
  MenuItem, 
  Box,
  ListItemIcon,
  ListItemText,
  Divider
} from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'
import CheckIcon from '@mui/icons-material/Check'
import LogoutIcon from '@mui/icons-material/Logout'
import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import { Anchor } from 'lucide-react'
import Home from './pages/Home'
import SignUp from './pages/SignUp'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import AddCertificate from './pages/AddCertificate'
import './App.css'
import { ColorModeContext } from './main'
import { supabase } from './supabaseClient'

interface UserData {
  id: string;
  forename: string;
  surname: string;
  email: string;
  nationality: string;
}

function App({ mode }: { mode: string }) {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [accountAnchorEl, setAccountAnchorEl] = useState<null | HTMLElement>(null)
  const [session, setSession] = useState<any>(null)
  const [userData, setUserData] = useState<UserData | null>(null)
  const colorMode = useContext(ColorModeContext)
  const navigate = useNavigate()

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session)
    })

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session)
    })

    return () => subscription.unsubscribe()
  }, [])

  useEffect(() => {
    const fetchUserData = async () => {
      if (session?.access_token) {
        try {
          const response = await fetch('/admin/users', {
            headers: {
              'Authorization': `Bearer ${session.access_token}`,
            },
          })
          if (response.ok) {
            const data = await response.json()
            // The API returns a list of users, but based on instructions 
            // "the user id will be retrieved automatically from the access token"
            // If the endpoint returns the current user directly or a list where we can find the user.
            // Usually such admin endpoints might return a list.
            // However, the prompt says "show the user name and email address".
            // If it returns a list, we might need to find the one matching session.user.id
            if (Array.isArray(data)) {
              const user = data.find(u => u.id === session.user.id)
              setUserData(user || null)
            } else {
              setUserData(data)
            }
          }
        } catch (error) {
          console.error('Error fetching user data:', error)
        }
      } else {
        setUserData(null)
      }
    }

    fetchUserData()
  }, [session])

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleAccountMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAccountAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
    setAccountAnchorEl(null)
  }

  const handleToggleDarkMode = () => {
    colorMode.toggleColorMode()
    handleClose()
  }

  const handleLogout = async () => {
    await supabase.auth.signOut()
    handleClose()
    navigate('/login')
  }

  return (
    <>
      <AppBar position="fixed" elevation={0} sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Toolbar>
          {session && (
            <IconButton
              size="large"
              edge="start"
              color="inherit"
              aria-label="menu"
              sx={{ mr: 2 }}
              onClick={handleMenu}
            >
              <MenuIcon />
            </IconButton>
          )}
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItem onClick={handleClose} component={RouterLink} to="/dashboard">
              Dashboard
            </MenuItem>
            <MenuItem onClick={handleClose} component={RouterLink} to="/add-certificate">
              Add Certificate
            </MenuItem>
            <MenuItem onClick={handleToggleDarkMode}>
              <ListItemIcon>
                {mode === 'dark' && <CheckIcon fontSize="small" />}
              </ListItemIcon>
              <ListItemText>Dark Mode</ListItemText>
            </MenuItem>
          </Menu>

          <Typography
            variant="h6"
            component={RouterLink}
            to="/"
            sx={{ 
              flexGrow: 1, 
              textDecoration: 'none', 
              color: 'inherit',
              display: 'flex',
              alignItems: 'center',
              gap: 1
            }}
          >
            <Anchor size={24} />
            SeaCert
          </Typography>

          <Box sx={{ display: 'flex', gap: 1 }}>
            {session ? (
              <>
                <Button
                  color="inherit"
                  onClick={handleAccountMenu}
                  startIcon={<AccountCircleIcon />}
                  sx={{ textTransform: 'none' }}
                >
                  My Account
                </Button>
                <Menu
                  anchorEl={accountAnchorEl}
                  open={Boolean(accountAnchorEl)}
                  onClose={handleClose}
                  anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'right',
                  }}
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                >
                  <Box sx={{ px: 2, py: 1, minWidth: 200 }}>
                    <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                      {userData ? `${userData.forename} ${userData.surname}` : 'User'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {userData?.email || session.user.email}
                    </Typography>
                  </Box>
                  <Divider />
                  <MenuItem onClick={handleLogout}>
                    <ListItemIcon>
                      <LogoutIcon fontSize="small" />
                    </ListItemIcon>
                    <ListItemText>Logout</ListItemText>
                  </MenuItem>
                </Menu>
              </>
            ) : (
              <>
                <Button 
                  color="inherit" 
                  component={RouterLink} 
                  to="/signup"
                >
                  Sign Up
                </Button>
                <Button 
                  color="secondary" 
                  variant="contained" 
                  component={RouterLink} 
                  to="/login"
                >
                  Login
                </Button>
              </>
            )}
          </Box>
        </Toolbar>
      </AppBar>
      <Toolbar /> {/* Spacer to prevent content from being hidden under fixed AppBar */}
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/login" element={<Login />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/add-certificate" element={<AddCertificate />} />
      </Routes>
    </>
  )
}

export default App
