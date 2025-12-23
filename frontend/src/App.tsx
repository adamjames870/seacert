import { useState, useEffect } from 'react'
import { Routes, Route, Link as RouterLink, useNavigate, Navigate } from 'react-router-dom'
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
  Divider,
  CircularProgress
} from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'
import LogoutIcon from '@mui/icons-material/Logout'
import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import EditIcon from '@mui/icons-material/Edit'
import { Anchor } from 'lucide-react'
import Home from './pages/Home'
import SignUp from './pages/SignUp'
import Login from './pages/Login'
import Certificates from './pages/Certificates.tsx'
import AddCertificate from './pages/AddCertificate'
import AddIssuer from './pages/AddIssuer'
import UpdateCertificate from './pages/UpdateCertificate'
import EditAccount from './pages/EditAccount'
import CertTypes from './pages/CertTypes'
import AddCertType from './pages/AddCertType'
import EditCertType from './pages/EditCertType'
import Issuers from './pages/Issuers'
import EditIssuer from './pages/EditIssuer'
import './App.css'
import { supabase } from './supabaseClient'
import { API_BASE_URL } from './config'

interface UserData {
  id: string;
  forename: string;
  surname: string;
  email: string;
  nationality: string;
  role?: string;
}

function App() {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [accountAnchorEl, setAccountAnchorEl] = useState<null | HTMLElement>(null)
  const [session, setSession] = useState<any>(undefined)
  const [userData, setUserData] = useState<UserData | null>(null)
  const [loadingUserData, setLoadingUserData] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session)
      if (!session) setLoadingUserData(false)
    })

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session)
      if (!session) {
        setUserData(null)
        setLoadingUserData(false)
      }
    })

    return () => subscription.unsubscribe()
  }, [])

  useEffect(() => {
    const fetchUserData = async () => {
      if (session?.access_token) {
        setLoadingUserData(true)
        try {
          const response = await fetch(`${API_BASE_URL}/admin/users`, {
            headers: {
              'Authorization': `Bearer ${session.access_token}`,
            },
          })
          if (response.ok) {
            const data = await response.json()
            if (Array.isArray(data)) {
              const user = data.find(u => u.id === session.user.id)
              setUserData(user || null)
            } else {
              setUserData(data)
            }
          }
        } catch (error) {
          console.error('Error fetching user data:', error)
        } finally {
          setLoadingUserData(false)
        }
      } else {
        setUserData(null)
        setLoadingUserData(false)
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

  const handleLogout = async () => {
    await supabase.auth.signOut()
    handleClose()
    navigate('/login')
  }

  const isAdmin = session?.user?.app_metadata?.role === 'admin'

  // Only block the whole app if we're waiting for the initial session check
  // or if we have a session but haven't started fetching user data yet.
  if (session === undefined || (session && loadingUserData && !userData)) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    )
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
            <MenuItem onClick={handleClose} component={RouterLink} to="/certificates">
              Certificates
            </MenuItem>
            <MenuItem onClick={handleClose} component={RouterLink} to="/add-certificate">
              Add Certificate
            </MenuItem>
            {isAdmin && (
              <>
                <Divider />
                <MenuItem onClick={handleClose} component={RouterLink} to="/cert-types">
                  Certificate Types
                </MenuItem>
                <MenuItem onClick={handleClose} component={RouterLink} to="/issuers">
                  Issuers
                </MenuItem>
              </>
            )}
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
                  <MenuItem onClick={handleClose} component={RouterLink} to="/edit-account">
                    <ListItemIcon>
                      <EditIcon fontSize="small" />
                    </ListItemIcon>
                    <ListItemText>Edit Account</ListItemText>
                  </MenuItem>
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
        <Route path="/" element={session ? <Navigate to="/certificates" replace /> : <Home />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/login" element={<Login />} />
        <Route path="/certificates" element={<Certificates />} />
        <Route path="/add-certificate" element={<AddCertificate />} />
        <Route path="/add-issuer" element={<AddIssuer />} />
        <Route path="/update-certificate/:id" element={<UpdateCertificate />} />
        
        {/* Admin only routes */}
        <Route path="/cert-types" element={isAdmin ? <CertTypes /> : <Navigate to="/certificates" replace />} />
        <Route path="/add-cert-type" element={isAdmin ? <AddCertType /> : <Navigate to="/certificates" replace />} />
        <Route path="/edit-cert-type/:id" element={isAdmin ? <EditCertType /> : <Navigate to="/certificates" replace />} />
        <Route path="/issuers" element={isAdmin ? <Issuers /> : <Navigate to="/certificates" replace />} />
        <Route path="/edit-issuer/:id" element={<EditIssuer />} />
        
        <Route path="/edit-account" element={<EditAccount />} />
      </Routes>
    </>
  )
}

export default App
