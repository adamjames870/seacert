import { useParams, Navigate, Link as RouterLink } from 'react-router-dom';
import { useEffect, useState, Suspense } from 'react';
import { Container, Box, Typography, Breadcrumbs, Link, CircularProgress, Divider } from '@mui/material';
import { ChevronRight, ArrowLeft } from 'lucide-react';

const BlogPost = () => {
  const { slug } = useParams();
  const [PostContent, setPostContent] = useState<any>(null);
  const [meta, setMeta] = useState<any>(null);
  const [notFound, setNotFound] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadPost = async () => {
      if (!slug) return;
      setLoading(true);

      try {
        // Use Vite's glob import to find all MDX files in the blog folder
        const modules = import.meta.glob('../blog/*.mdx');
        let modulePath = '';

        // Find the module that matches the slug in its meta
        for (const path in modules) {
          const module: any = await modules[path]();
          if (module.meta?.slug === slug) {
            modulePath = path;
            setMeta(module.meta);
            setPostContent(() => module.default);
            
            document.title = `${module.meta.title} | SeaCert`;
            const metaDesc = document.querySelector('meta[name="description"]');
            if (metaDesc) {
              metaDesc.setAttribute('content', module.meta.description);
            }
            break;
          }
        }

        if (!modulePath) {
          setNotFound(true);
        }
      } catch (err) {
        console.error('Error loading blog post:', err);
        setNotFound(true);
      } finally {
        setLoading(false);
      }
    };

    loadPost();
  }, [slug]);

  if (notFound) {
    return <Navigate to="/blog" replace />;
  }

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container maxWidth="md" sx={{ mt: 4, mb: 10 }}>
      <Box sx={{ mb: 4 }}>
        <Breadcrumbs 
          separator={<ChevronRight size={16} />} 
          aria-label="breadcrumb"
          sx={{ mb: 3 }}
        >
          <Link component={RouterLink} underline="hover" color="inherit" to="/">
            Home
          </Link>
          <Link component={RouterLink} underline="hover" color="inherit" to="/blog">
            Blog
          </Link>
          <Typography color="text.primary">{meta?.title || 'Article'}</Typography>
        </Breadcrumbs>

        <Link 
          component={RouterLink} 
          to="/blog" 
          sx={{ 
            display: 'inline-flex', 
            alignItems: 'center', 
            gap: 1, 
            mb: 4, 
            textDecoration: 'none',
            color: 'text.secondary',
            '&:hover': { color: 'primary.main' }
          }}
        >
          <ArrowLeft size={18} />
          Back to Blog
        </Link>
      </Box>

      <Box component="article">
        <Suspense fallback={
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
            <CircularProgress />
          </Box>
        }>
          <Box className="blog-content" sx={{ 
            '& h1': { variant: 'h3', fontWeight: 800, mb: 4 },
            '& h2': { variant: 'h4', fontWeight: 700, mt: 6, mb: 3 },
            '& h3': { variant: 'h5', fontWeight: 600, mt: 4, mb: 2 },
            '& p': { variant: 'body1', mb: 2, lineHeight: 1.7 },
            '& ul, & ol': { mb: 3, pl: 4 },
            '& li': { mb: 1 },
            '& hr': { my: 6, opacity: 0.1 }
          }}>
            <PostContent />
          </Box>
        </Suspense>
      </Box>
      
      <Divider sx={{ my: 8 }} />
      
      <Box sx={{ textAlign: 'center' }}>
        <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
          Found this guide helpful?
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          Share it with your crewmates or colleagues.
        </Typography>
      </Box>
    </Container>
  );
};

export default BlogPost;
