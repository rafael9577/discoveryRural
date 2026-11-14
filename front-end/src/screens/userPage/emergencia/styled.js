import { StyleSheet } from 'react-native';

export const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: '#F5F7FA'
  },
  
  header: { 
    flexDirection: 'row', 
    alignItems: 'center', 
    justifyContent: 'center',
    paddingHorizontal: 16,
    paddingVertical: 15,
    backgroundColor: 'transparent',
    borderBottomWidth: 1,
    borderBottomColor: '#E0E4E8',
  },
  headerTitleBox: { 
    alignItems: 'center',
    justifyContent: 'center',
  },
  title: { 
    fontSize: 20,
    fontWeight: 'bold',
    color: '#2C3E50',
    textTransform: 'uppercase',
  },
  subtitle: { 
    fontSize: 13, 
    color: '#7F8C8D', 
    marginTop: 4,
  },

  scrollContent: { 
    padding: 16, 
    paddingBottom: 40,
  },

  section: { 
    marginBottom: 24 
  },
  sectionTitle: { 
    fontSize: 16, 
    fontWeight: 'bold', 
    color: '#E74C3C',
    marginBottom: 12, 
    borderBottomWidth: 2, 
    borderBottomColor: '#E74C3C', 
    paddingBottom: 6,
    textTransform: 'uppercase',
  },

  /* CARDS */
  card: { 
    backgroundColor: '#FFFFFF', 
    borderRadius: 16,
    padding: 16, 
    marginBottom: 12, 
    shadowColor: '#000', 
    shadowOffset: { width: 0, height: 2 }, 
    shadowOpacity: 0.1, 
    shadowRadius: 4, 
    elevation: 3,
    borderWidth: 1,
    borderColor: '#EAEAEA',
  },
  cardTitle: { 
    fontSize: 16, 
    fontWeight: 'bold', 
    color: '#34495E', 
    marginBottom: 10 
  },
  cardText: { 
    fontSize: 14, 
    color: '#555555', 
    lineHeight: 22,
    marginBottom: 6 
  },
  highlightText: { 
    fontWeight: 'bold', 
    color: '#D35400'
  },
});