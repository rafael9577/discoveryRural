import { StyleSheet } from 'react-native';

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F5F5F5',
  },
  
  /* HEADER STYLES */
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: '#FFF',
    borderBottomWidth: 1,
    borderBottomColor: '#E0E0E0',
    zIndex: 2, // Garante que o header fique acima do mapa
  },
  iconButton: {
    padding: 8,
  },
  headerTitleBox: {
    borderWidth: 1,
    borderColor: '#000',
    paddingVertical: 8,
    paddingHorizontal: 32,
    backgroundColor: '#FFF',
  },
  headerTitleText: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#000',
  },
  menuButton: {
    padding: 8,
    borderWidth: 1,
    borderColor: '#000',
    borderRadius: 4,
    backgroundColor: '#FFF',
  },
  menuButtonAtivo: {
    backgroundColor: '#228B22', 
    borderColor: '#228B22',
  },

  /* MAP STYLES */
  mapContainer: {
    flex: 1,
    margin: 16,
    borderRadius: 16,
    overflow: 'hidden', 
    borderWidth: 1,
    borderColor: '#000',
  },
  map: {
    width: '100%',
    height: '100%',
  },

  /* OVERLAY MENU STYLES */
  overlayMenu: {
    position: 'absolute',
    top: 0,
    right: 0,
    width: '75%', 
    backgroundColor: '#FFF',
    borderBottomLeftRadius: 8,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
    elevation: 5,
  },
  overlayHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: '#228B22', 
    paddingHorizontal: 16,
    paddingVertical: 12,
  },
  overlayHeaderText: {
    color: '#FFF',
    fontWeight: 'bold',
    fontSize: 14,
  },
  overlayBody: {
    backgroundColor: '#EAEAEA', 
  },
  filterTitleContainer: {
    backgroundColor: '#D3D3D3',
    paddingVertical: 8,
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#B0B0B0',
  },
  filterTitleText: {
    fontWeight: 'bold',
    fontSize: 12,
    color: '#000',
  },
  checkboxItemContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 12,
    paddingHorizontal: 16,
  },
  checkboxItemText: {
    fontSize: 14,
    color: '#000',
    fontWeight: '500',
  },
  checkboxBox: {
    width: 20,
    height: 20,
    borderWidth: 1,
    borderColor: '#000',
    backgroundColor: '#FFF',
    justifyContent: 'center',
    alignItems: 'center',
  },
  separator: {
    height: 1,
    backgroundColor: '#999',
    marginHorizontal: 16,
  },
});