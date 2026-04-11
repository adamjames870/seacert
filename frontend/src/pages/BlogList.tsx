import { useState, useEffect } from 'react';
import { Container, Typography, Box, Grid, Card, CardContent, CardActionArea, Divider, CircularProgress } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { BookOpen, Calendar } from 'lucide-react';

const BlogList = () => {
  const [posts, setPosts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadPosts = async () => {
      // Use Vite's glob import to find all MDX files in the blog folder
      const modules = import.meta.glob('../blog/*.mdx');
      const postList = [];

      for (const path in modules) {
        const module: any = await modules[path]();
        if (module.meta) {
          postList.push(module.meta);
        }
      }

      setPosts(postList);
      setLoading(false);
    };

    loadPosts();
  }, []);

  if (loading) {
    return (
      <Container maxWidth="md" sx={{ mt: 10, textAlign: 'center' }}>
        <CircularProgress />
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ mt: 6, mb: 10 }}>
      <Box sx={{ mb: 6, textAlign: 'center' }}>
        <Typography variant="h3" component="h1" gutterBottom sx={{ fontWeight: 800 }}>
          SeaCert Blog
        </Typography>
        <Typography variant="h6" color="text.secondary">
          Insights, guides, and updates for seafarers.
        </Typography>
      </Box>

      <Grid container spacing={4}>
        {posts.map((post) => (
          <Grid size={12} key={post.slug}>
            <Card sx={{ 
              borderRadius: 2, 
              boxShadow: '0 4px 20px rgba(0,0,0,0.08)',
              transition: 'transform 0.2s',
              '&:hover': { transform: 'translateY(-4px)' }
            }}>
              <CardActionArea component={RouterLink} to={`/blog/${post.slug}`}>
                <CardContent sx={{ p: 4 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2, color: 'primary.main' }}>
                    <BookOpen size={18} />
                    <Typography variant="overline" sx={{ fontWeight: 700 }}>Guide</Typography>
                  </Box>
                  <Typography variant="h4" component="h2" gutterBottom sx={{ fontWeight: 700 }}>
                    {post.title}
                  </Typography>
                  <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
                    {post.description}
                  </Typography>
                  <Divider sx={{ mb: 2 }} />
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, color: 'text.disabled' }}>
                    <Calendar size={14} />
                    <Typography variant="caption">April 9, 2026</Typography>
                  </Box>
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Container>
  );
};

export default BlogList;
