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
  CircularProgress,
  Container,
  Link
} from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'
import LogoutIcon from '@mui/icons-material/Logout'
import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import EditIcon from '@mui/icons-material/Edit'
import PeopleIcon from '@mui/icons-material/People'
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
import AdminUsers from './pages/AdminUsers'
import AdminUserCertificates from './pages/AdminUserCertificates'
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
        console.log('Fetching user data for session:', session.user.id);
        try {
          const response = await fetch(`${API_BASE_URL}/admin/users`, {
            headers: {
              'Authorization': `Bearer ${session.access_token}`,
            },
          })
          console.log('User data fetch response status:', response.status);
          if (response.ok) {
            const data = await response.json()
            console.log('User data fetch response data:', data);
            
            let usersArray: UserData[] = [];
            if (Array.isArray(data)) {
              usersArray = data;
            } else if (data && typeof data === 'object') {
              // Try to find an array in the object properties
              const arrayKey = Object.keys(data).find(key => Array.isArray(data[key]));
              if (arrayKey) {
                usersArray = data[arrayKey];
              } else if (data.id || data.email) {
                // If it's a single user object
                setUserData(data as UserData);
                return;
              }
            }

            if (usersArray.length > 0) {
              const user = usersArray.find(u => u.id === session.user.id)
              if (user) {
                console.log('Found user in /admin/users list:', user);
                setUserData(user)
              } else {
                console.warn('User session ID not found in /admin/users list');
                // Fallback: If we got some data but not the specific user, 
                // and it was a single object response, we already set it above and returned.
                // If it was an array but our ID wasn't in it, something is wrong.
                setUserData(null)
              }
            } else if (userData) {
              // already set via single object check
              return;
            } else {
              console.warn('No user data found in /admin/users response');
              setUserData(null)
            }
          } else {
            console.warn('Failed to fetch user data via /admin/users - status:', response.status);
            // Try fetching from public user endpoint if it exists or just settle for session data
            setUserData(null)
          }
        } catch (error) {
          console.error('Error fetching user data:', error)
          setUserData(null)
        } finally {
          setLoadingUserData(false)
        }
      } else {
        console.log('No session, clearing user data.');
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

  const isAdmin = session?.user?.app_metadata?.role === 'admin' || userData?.role === 'admin' || session?.user?.user_metadata?.role === 'admin'
  console.log('Admin status:', { 
    isAdmin, 
    app_metadata_role: session?.user?.app_metadata?.role, 
    user_metadata_role: session?.user?.user_metadata?.role,
    userData_role: userData?.role 
  });

  // Only block the whole app if we're waiting for the initial session check
  if (session === undefined) {
    console.log('App loading: session undefined');
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    )
  }

  const isAuthPath = ['/login', '/signup', '/'].includes(window.location.pathname);

  // If we have a session but we're still loading user data and don't have it yet, show loading spinner.
  // We'll also allow it to proceed if loadingUserData becomes false, regardless of if userData is found.
  if (session && loadingUserData && !userData && !isAuthPath) {
    console.log('App loading: session exists, fetching user data...');
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
            {isAdmin && [
              <Divider key="divider" />,
              <MenuItem key="issuers" onClick={handleClose} component={RouterLink} to="/issuers">
                Issuers
              </MenuItem>,
              <MenuItem key="cert-types" onClick={handleClose} component={RouterLink} to="/cert-types">
                Certificate Types
              </MenuItem>,
              <MenuItem key="users" onClick={handleClose} component={RouterLink} to="/admin/users">
                <ListItemIcon>
                  <PeopleIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText>User Management</ListItemText>
              </MenuItem>
            ]}
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
              gap: 1,
              fontSize: { xs: '1.1rem', sm: '1.25rem' }
            }}
          >
            <Anchor size={20} />
            <Box component="span" sx={{ display: { xs: 'none', sm: 'inline' }, fontWeight: 700 }}>SeaCert</Box>
            <Box component="span" sx={{ display: { xs: 'inline', sm: 'none' }, fontWeight: 700 }}>SC</Box>
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
                  <Box component="span" sx={{ display: { xs: 'none', sm: 'inline' } }}>My Account</Box>
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
        
        {/* Certificate Types: Admin can see all, Users can add types from AddCertificate flow */}
        <Route path="/cert-types" element={isAdmin ? <CertTypes /> : <Navigate to="/certificates" replace />} />
        <Route path="/add-cert-type" element={session ? <AddCertType /> : <Navigate to="/login" replace />} />
        <Route path="/edit-cert-type/:id" element={isAdmin ? <EditCertType /> : <Navigate to="/certificates" replace />} />
        <Route path="/issuers" element={isAdmin ? <Issuers /> : <Navigate to="/certificates" replace />} />
        <Route path="/edit-issuer/:id" element={<EditIssuer />} />
        
        <Route path="/admin/users" element={isAdmin ? <AdminUsers /> : <Navigate to="/certificates" replace />} />
        <Route path="/admin/users/:userId/certificates" element={isAdmin ? <AdminUserCertificates /> : <Navigate to="/certificates" replace />} />

        <Route path="/edit-account" element={<EditAccount />} />
      </Routes>
      
      <Box component="footer" sx={{ py: 3, borderTop: 1, borderColor: 'divider', mt: 'auto', bgcolor: 'background.paper' }}>
        <Container maxWidth="lg">
          <Typography variant="body2" color="text.secondary" align="center">
            © {new Date().getFullYear()} SeaCert. Contact us: {' '}
            <Link href="mailto:hello@seacert.app" color="inherit" sx={{ fontWeight: 600 }}>
              hello@seacert.app
            </Link>
          </Typography>
        </Container>
      </Box>
    </>
  )
}

export default App
