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
import DashboardIcon from '@mui/icons-material/Dashboard'
import CardMembershipIcon from '@mui/icons-material/CardMembership'
import AddCircleIcon from '@mui/icons-material/AddCircle'
import AssessmentIcon from '@mui/icons-material/Assessment'
import HistoryIcon from '@mui/icons-material/History'
import AddLocationIcon from '@mui/icons-material/AddLocation'
import ArticleIcon from '@mui/icons-material/Article'
import BusinessIcon from '@mui/icons-material/Business'
import CategoryIcon from '@mui/icons-material/Category'
import DirectionsBoatIcon from '@mui/icons-material/DirectionsBoat'
import SearchIcon from '@mui/icons-material/Search'
import { Anchor, Shield, Ship as ShipIcon } from 'lucide-react'
import Home from './pages/Home'
import Dashboard from './pages/Dashboard'
import SignUp from './pages/SignUp'
import Login from './pages/Login'
import Certificates from './pages/Certificates'
import AddCertificate from './pages/AddCertificate'
import AddIssuer from './pages/AddIssuer'
import UpdateCertificate from './pages/UpdateCertificate'
import EditAccount from './pages/EditAccount'
import CertTypes from './pages/CertTypes'
import AddCertType from './pages/AddCertType'
import EditCertType from './pages/EditCertType'
import Issuers from './pages/Issuers'
import EditIssuer from './pages/EditIssuer'
import SeatimeHistory from './pages/SeatimeHistory'
import AddSeatime from './pages/AddSeatime'
import UpdateSeatime from './pages/UpdateSeatime'
import ManageSeatimeLookups from './pages/ManageSeatimeLookups'
import Ships from './pages/Ships'
import ShipForm from './pages/ShipForm'
import Privacy from './pages/Privacy'
import Terms from './pages/Terms'
import ForgotPassword from './pages/ForgotPassword'
import ResetPassword from './pages/ResetPassword'
import BlogList from './pages/BlogList'
import BlogPost from './pages/BlogPost'
import CertificateWizard from './pages/wizard/CertificateWizard'
import ReportPreviewDialog from './components/ReportPreviewDialog'
import CookieConsent from './components/CookieConsent'
import EmailConsentModal from './components/EmailConsentModal'
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
  email_consent?: boolean;
  email_consent_timestamp?: string;
  email_consent_version?: string;
  email_consent_source?: string;
}

function App() {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [accountAnchorEl, setAccountAnchorEl] = useState<null | HTMLElement>(null)
  const [reportDialogOpen, setReportDialogOpen] = useState(false)
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

  const handleOpenReport = () => {
    setReportDialogOpen(true)
    handleClose()
  }

  const handleLogout = async () => {
    try {
      const { error } = await supabase.auth.signOut()
      if (error) {
        console.error('Error during signOut:', error.message)
      }
    } catch (err) {
      console.error('Unexpected error during logout:', err)
    } finally {
      // Clear all local state regardless of server response
      setSession(null)
      setUserData(null)
      
      // Explicitly clear any Supabase related items from localStorage
      // This is a "bulletproof" step to ensure no stale session remains
      Object.keys(localStorage).forEach(key => {
        if (key.startsWith('sb-')) {
          localStorage.removeItem(key)
        }
      })
      
      handleClose()
      navigate('/')
    }
  }

  const isAdmin = session?.user?.app_metadata?.role === 'admin'

  const handleConsentClose = (updatedUser?: UserData) => {
    if (updatedUser) {
      setUserData(updatedUser);
    } else {
      // Re-fetch user data if for some reason modal didn't pass it back
      const reFetch = async () => {
        if (session?.access_token) {
          try {
            const response = await fetch(`${API_BASE_URL}/admin/users`, {
              headers: {
                'Authorization': `Bearer ${session.access_token}`,
              },
            });
            if (response.ok) {
              const data = await response.json();
              if (Array.isArray(data)) {
                const user = data.find(u => u.id === session.user.id);
                setUserData(user || null);
              } else {
                setUserData(data);
              }
            }
          } catch (e) {
            console.error(e);
          }
        }
      };
      reFetch();
    }
  }

  // Only block the whole app if we're waiting for the initial session check
  // or if we have a session but haven't started fetching user data yet.
  if (session === undefined || (session && loadingUserData && !userData)) {
    return (
      <>
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
          <CircularProgress />
        </Box>
        <CookieConsent />
        <ReportPreviewDialog open={reportDialogOpen} onClose={() => setReportDialogOpen(false)} />
      </>
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
            disableScrollLock={true}
          >
            <MenuItem onClick={handleClose} component={RouterLink} to="/">
              <ListItemIcon>
                <DashboardIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Dashboard</ListItemText>
            </MenuItem>

            <Divider />

            <MenuItem onClick={handleClose} component={RouterLink} to="/certificates">
              <ListItemIcon>
                <CardMembershipIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Certificates</ListItemText>
            </MenuItem>
            <MenuItem onClick={handleClose} component={RouterLink} to="/certificate-wizard">
              <ListItemIcon>
                <AddCircleIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Add Certificate</ListItemText>
            </MenuItem>
            <MenuItem onClick={handleOpenReport}>
              <ListItemIcon>
                <AssessmentIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Certificate Report</ListItemText>
            </MenuItem>
            
            <Divider />
            
            <MenuItem onClick={handleClose} component={RouterLink} to="/seatime">
              <ListItemIcon>
                <HistoryIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Seatime History</ListItemText>
            </MenuItem>
            <MenuItem onClick={handleClose} component={RouterLink} to="/add-seatime">
              <ListItemIcon>
                <AddLocationIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Add Seatime Period</ListItemText>
            </MenuItem>
            
            <Divider />

            <MenuItem onClick={handleClose} component={RouterLink} to="/blog">
              <ListItemIcon>
                <ArticleIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Blog</ListItemText>
            </MenuItem>
            
            {isAdmin && [
              <Divider key="divider-admin" />,
              <Box key="admin-certs-header" sx={{ px: 2, py: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                <Shield size={16} color="gray" />
                <Typography variant="overline" sx={{ fontWeight: 700, color: 'text.secondary' }}>Admin: Certificates</Typography>
              </Box>,
              <MenuItem key="issuers" onClick={handleClose} component={RouterLink} to="/issuers" sx={{ pl: 4 }}>
                <ListItemIcon>
                  <BusinessIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText>Issuers</ListItemText>
              </MenuItem>,
              <MenuItem key="cert-types" onClick={handleClose} component={RouterLink} to="/cert-types" sx={{ pl: 4 }}>
                <ListItemIcon>
                  <CategoryIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText>Certificate Types</ListItemText>
              </MenuItem>,
              
              <Box key="admin-seatime-header" sx={{ px: 2, py: 1, mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                <ShipIcon size={16} color="gray" />
                <Typography variant="overline" sx={{ fontWeight: 700, color: 'text.secondary' }}>Admin: Seatime & Ships</Typography>
              </Box>,
              <MenuItem key="ships" onClick={handleClose} component={RouterLink} to="/ships" sx={{ pl: 4 }}>
                <ListItemIcon>
                  <DirectionsBoatIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText>Ships</ListItemText>
              </MenuItem>,
              <MenuItem key="seatime-lookups" onClick={handleClose} component={RouterLink} to="/admin/seatime-lookups" sx={{ pl: 4 }}>
                <ListItemIcon>
                  <SearchIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText>Seatime Lookups</ListItemText>
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
            {!session && (
              <Button 
                color="inherit" 
                component={RouterLink} 
                to="/blog"
                sx={{ display: { xs: 'none', md: 'inline-flex' }, mr: 1 }}
              >
                Blog
              </Button>
            )}
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
                  disableScrollLock={true}
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
        <Route path="/" element={session ? <Dashboard /> : <Home />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/login" element={<Login />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/reset-password" element={<ResetPassword />} />
        <Route path="/certificates" element={<Certificates />} />
        <Route path="/add-certificate" element={<AddCertificate />} />
        <Route path="/certificate-wizard" element={session ? <CertificateWizard /> : <Navigate to="/login" replace />} />
        <Route path="/add-issuer" element={<AddIssuer />} />
        <Route path="/update-certificate/:id" element={<UpdateCertificate />} />

        <Route path="/seatime" element={session ? <SeatimeHistory /> : <Navigate to="/login" replace />} />
        <Route path="/add-seatime" element={session ? <AddSeatime /> : <Navigate to="/login" replace />} />
        <Route path="/update-seatime/:id" element={session ? <UpdateSeatime /> : <Navigate to="/login" replace />} />
        <Route path="/ships" element={isAdmin ? <Ships /> : <Navigate to="/certificates" replace />} />
        <Route path="/add-ship" element={isAdmin ? <ShipForm /> : <Navigate to="/certificates" replace />} />
        <Route path="/edit-ship/:id" element={isAdmin ? <ShipForm /> : <Navigate to="/certificates" replace />} />
        <Route path="/admin/seatime-lookups" element={isAdmin ? <ManageSeatimeLookups /> : <Navigate to="/certificates" replace />} />
        
        {/* Certificate Types: Admin can see all, Users can add types from AddCertificate flow */}
        <Route path="/cert-types" element={isAdmin ? <CertTypes /> : <Navigate to="/certificates" replace />} />
        <Route path="/add-cert-type" element={session ? <AddCertType /> : <Navigate to="/login" replace />} />
        <Route path="/edit-cert-type/:id" element={isAdmin ? <EditCertType /> : <Navigate to="/certificates" replace />} />
        <Route path="/issuers" element={isAdmin ? <Issuers /> : <Navigate to="/certificates" replace />} />
        <Route path="/edit-issuer/:id" element={<EditIssuer />} />
        
        <Route path="/edit-account" element={<EditAccount />} />
        <Route path="/blog" element={<BlogList />} />
        <Route path="/blog/:slug" element={<BlogPost />} />
        <Route path="/privacy" element={<Privacy />} />
        <Route path="/terms" element={<Terms />} />
      </Routes>
      
      <Box component="footer" sx={{ py: 3, borderTop: 1, borderColor: 'divider', mt: 'auto', bgcolor: 'background.paper' }}>
        <Container maxWidth="lg">
          <Typography variant="body2" color="text.secondary" align="center">
            © {new Date().getFullYear()} SeaCert. {' '}
            <Link component={RouterLink} to="/privacy" color="inherit" sx={{ fontWeight: 600, mr: 2 }}>
              Privacy Policy
            </Link>
            <Link component={RouterLink} to="/terms" color="inherit" sx={{ fontWeight: 600, mr: 2 }}>
              Terms & Conditions
            </Link>
            Contact us: {' '}
            <Link href="mailto:hello@seacert.app" color="inherit" sx={{ fontWeight: 600 }}>
              hello@seacert.app
            </Link>
          </Typography>
        </Container>
      </Box>
      <CookieConsent />
      <EmailConsentModal 
        open={!!(session && userData && userData.email_consent_timestamp === null)} 
        onClose={handleConsentClose} 
      />
      <ReportPreviewDialog open={reportDialogOpen} onClose={() => setReportDialogOpen(false)} />
    </>
  )
}

export default App
