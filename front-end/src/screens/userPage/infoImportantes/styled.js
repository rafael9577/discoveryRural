import { StyleSheet } from 'react-native';

export const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: '#F5F7FA' 
  },
  scrollContent: { 
    padding: 20, 
    paddingBottom: 40 
  },
  header: { 
    flexDirection: 'row', 
    alignItems: 'center', 
    marginBottom: 20, 
    marginTop: 10 
  },
  backButton: 
  { padding: 10, 
    marginRight: 10 
  },
  headerTextContainer: { 
    flex: 1 
  },
  title: { fontSize: 24
    , fontWeight: 'bold'
    , color: '#2C3E50' 
  },
  subtitle: { 
    fontSize: 14, 
    color: '#7F8C8D', 
    marginTop: 4 
  },
  section: { 
    marginBottom: 24 
  },
  sectionTitle: { 
    fontSize: 18, 
    fontWeight: 'bold', 
    color: '#E74C3C', 
    marginBottom: 12, 
    borderBottomWidth: 1, 
    borderBottomColor: '#E74C3C', 
    paddingBottom: 4 
  },
  card: { 
    backgroundColor: '#FFFFFF', 
    borderRadius: 12, 
    padding: 16, 
    marginBottom: 12, 
    shadowColor: '#000', 
    shadowOffset: { width: 0, height: 2 }, 
    shadowOpacity: 0.1, 
    shadowRadius: 4, 
    elevation: 3 
  },
  cardTitle: { 
    fontSize: 16, 
    fontWeight: 'bold', 
    color: '#34495E', 
    marginBottom: 8 
  },
  cardText: { fontSize: 14, 
    color: '#555555', 
    lineHeight: 20, 
    marginBottom: 4 
  },
  highlightText: { 
    fontWeight: 'bold', 
    color: '#D35400' 
  },
  footer: { alignItems: 'center', 
    marginTop: 20, 
    paddingTop: 20, 
    borderTopWidth: 1, 
    borderTopColor: '#BDC3C7' 
  },
  footerText: { 
    fontSize: 12, 
    color: '#95A5A6', 
    textAlign: 'center', 
    marginBottom: 16, 
    lineHeight: 18 
  },
  logo: { 
    width: 150, 
    height: 50 
  },
});